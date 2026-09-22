// SPDX-FileCopyrightText: 2022-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package metricdata

import (
	"fmt"

	"github.com/omec-project/metricfunc/internal/promclient"
	"github.com/omec-project/metricfunc/logger"
	"github.com/omec-project/util/metricinfo"
)

// smSessKey identifies one smf_pdu_sessions Prometheus series by its label values.
type smSessKey struct {
	smfIP, slice, dnn, upf string
}

// smSessCounts holds the live session count per smf_pdu_sessions label combination.
// A single shared counter can't back a multi-label gauge: stamping one process-wide
// total onto whichever tuple triggered the last event leaves every other tuple's
// series frozen at a stale value and never reporting the deregistration that
// actually happened. Access is guarded by metricData.SubLock, which every caller
// below already holds.
var smSessCounts = map[smSessKey]uint64{}

func smSessKeyOf(sub *metricinfo.CoreSubscriber) smSessKey {
	return smSessKey{smfIP: sub.SmfIp, slice: sub.Slice, dnn: sub.Dnn, upf: sub.UpfName}
}

func incSmSessCount(sub *metricinfo.CoreSubscriber) {
	key := smSessKeyOf(sub)
	smSessCounts[key]++
	promclient.SetSmfSessStats(key.smfIP, key.slice, key.dnn, key.upf, smSessCounts[key])
}

// decSmSessCount lowers the count for sub's label tuple and drops the series entirely
// once it reaches zero, rather than leaving it published at a count that no longer
// reflects any session.
func decSmSessCount(sub *metricinfo.CoreSubscriber) {
	key := smSessKeyOf(sub)
	if smSessCounts[key] == 0 {
		return
	}
	smSessCounts[key]--
	if smSessCounts[key] == 0 {
		delete(smSessCounts, key)
		promclient.DeleteSmfSessStats(key.smfIP, key.slice, key.dnn, key.upf)
		return
	}
	promclient.SetSmfSessStats(key.smfIP, key.slice, key.dnn, key.upf, smSessCounts[key])
}

func HandleSubscriberEvent(subsData *metricinfo.CoreSubscriberData, sourceNf metricinfo.NfType) {
	// An imsi is this store's key. An event without one can't be merged into any real
	// subscriber's record and, left in the map, sits there forever as a row no later event can
	// ever reach again - so it can only be discarded here, not stored.
	if subsData.Subscriber.Imsi == "" {
		logger.CacheLog.Warnf("dropping subscriber event with empty imsi from sourceNF [%v]", sourceNf)
		return
	}

	switch subsData.Operation {
	case metricinfo.SubsOpAdd:
		storeSubscriber(&subsData.Subscriber, sourceNf)
	case metricinfo.SubsOpMod:
		updateSubscriber(&subsData.Subscriber, sourceNf)
	case metricinfo.SubsOpDel:
		err := deleteSubscriber(&subsData.Subscriber)
		if err != nil {
			logger.CacheLog.Infof("delete subscriber %v failed for sourceNF [%v]", subsData.Subscriber.Imsi, sourceNf)
		}
	default:
		logger.CacheLog.Errorf("unknown smf subscriber operation [%v]", subsData.Operation)
	}
}

func storeSubscriber(sub *metricinfo.CoreSubscriber, sourceNf metricinfo.NfType) {
	metricData.SubLock.Lock()

	if _, ok := metricData.Subscribers[sub.Imsi]; !ok {
		addSubscriberLocked(sub, sourceNf)
		metricData.SubLock.Unlock()
	} else {
		metricData.SubLock.Unlock()
		updateSubscriber(sub, sourceNf)
	}
}

// addSubscriberLocked stores a brand-new subscriber entry. Caller must hold metricData.SubLock.
func addSubscriberLocked(sub *metricinfo.CoreSubscriber, sourceNf metricinfo.NfType) {
	metricData.Subscribers[sub.Imsi] = sub

	// Only an SMF-sourced event carries real PDU session context
	if sourceNf == metricinfo.NfTypeSmf {
		incSmSessCount(sub)
	}
	logger.CacheLog.Debugf("storing subscriber with imsi [%s]", sub.Imsi)
	pushPrometheusCoreSubData(sub)
}

func updateSubscriber(sub *metricinfo.CoreSubscriber, sourceNf metricinfo.NfType) {
	metricData.SubLock.Lock()
	defer metricData.SubLock.Unlock()
	s, ok := metricData.Subscribers[sub.Imsi]
	if !ok {
		// A Mod can legitimately arrive before any Add for this imsi (e.g. AMF only learns
		// SUPI once authentication completes, so its first published event for a subscriber
		// is often a Mod). Upsert instead of dropping it, or this subscriber's data never
		// makes it into core_subscriber at all.
		addSubscriberLocked(sub, sourceNf)
		return
	}

	deletePrometheusCoreSubData(s)

	switch sourceNf {
	case metricinfo.NfTypeSmf:
		// SMF specific fields
		oldKey := smSessKeyOf(s)
		fillSmfSubsriberData(sub, s)
		if newKey := smSessKeyOf(s); newKey != oldKey {
			decSmSessCount(&metricinfo.CoreSubscriber{SmfIp: oldKey.smfIP, Slice: oldKey.slice, Dnn: oldKey.dnn, UpfName: oldKey.upf})
			incSmSessCount(s)
		}
	case metricinfo.NfTypeAmf:
		// AMF specific fields
		fillAmfSubsriberData(sub, s)
	}
	pushPrometheusCoreSubData(s)
}

func deleteSubscriber(sub *metricinfo.CoreSubscriber) error {
	metricData.SubLock.Lock()
	defer metricData.SubLock.Unlock()
	imsi := sub.Imsi
	s, ok := metricData.Subscribers[imsi]
	if !ok {
		return fmt.Errorf("subscriber with imsi [%s] already deleted", imsi)
	}

	decSmSessCount(s)
	deletePrometheusCoreSubData(s)
	s.SmfSubState = sub.SmfSubState
	s.AmfSubState = sub.AmfSubState

	// register disconnect state
	pushPrometheusCoreSubData(s)
	delete(metricData.Subscribers, imsi)

	// register subscriber delete
	deletePrometheusCoreSubData(s)

	logger.CacheLog.Debugf("deleting subscriber with imsi [%s]", imsi)

	return nil
}

func GetSubscriber(key string) (*metricinfo.CoreSubscriber, error) {
	metricData.SubLock.RLock()
	defer metricData.SubLock.RUnlock()
	if sub, ok := metricData.Subscribers[key]; ok {
		return sub, nil
	}
	return nil, fmt.Errorf("subscriber with key [%v] not found", key)
}

func GetSubscriberImsiFromIpAddr(ipaddr string) (*metricinfo.CoreSubscriber, error) {
	metricData.SubLock.RLock()
	defer metricData.SubLock.RUnlock()
	for imsi, sub := range metricData.Subscribers {
		if sub.IPAddress == ipaddr {
			logger.CacheLog.Infof("found subscriber with ip-addr [%s], imsi [%s]", ipaddr, imsi)
			return sub, nil
		}
	}
	return nil, fmt.Errorf("subscriber with ip-addr [%v] not found", ipaddr)
}

func GetSubscriberAll() []string {
	imsis := []string{}
	metricData.SubLock.RLock()
	defer metricData.SubLock.RUnlock()

	for imsi := range metricData.Subscribers {
		imsis = append(imsis, imsi)
	}

	return imsis
}

// Pushing to prometheus client module
func pushPrometheusCoreSubData(sub *metricinfo.CoreSubscriber) {
	promclient.PushCoreSubData(sub.Imsi, sub.IPAddress, sub.SmfSubState, sub.SmfIp, sub.Dnn, sub.Slice, sub.UpfName)
}

// Pushing to prometheus client module
func deletePrometheusCoreSubData(sub *metricinfo.CoreSubscriber) {
	promclient.DeleteCoreSubData(sub.Imsi, sub.IPAddress, sub.SmfSubState, sub.SmfIp, sub.Dnn, sub.Slice, sub.UpfName)
}

func fillSmfSubsriberData(s, d *metricinfo.CoreSubscriber) {
	// ip-addr
	if s.IPAddress != "" {
		d.IPAddress = s.IPAddress
	}

	// slice
	if s.Slice != "" {
		d.Slice = s.Slice
	}

	// dnn
	if s.Dnn != "" {
		d.Dnn = s.Dnn
	}

	// upf name
	if s.UpfName != "" {
		d.UpfName = s.UpfName
	}

	// upf ip
	if s.UpfAddr != "" {
		d.UpfAddr = s.UpfAddr
	}

	// always overwrite subscriber state
	d.SmfSubState = s.SmfSubState
}

func fillAmfSubsriberData(s, d *metricinfo.CoreSubscriber) {
	// AmfId
	if s.AmfId != "" {
		d.AmfId = s.AmfId
	}

	// Guti
	if s.Guti != "" {
		d.Guti = s.Guti
	}

	// TMSI
	if s.Tmsi != 0 {
		d.Tmsi = s.Tmsi
	}

	// Amf Ngap Id
	if s.AmfNgapId != 0 {
		d.AmfNgapId = s.AmfNgapId
	}

	// Ran Ngap Id
	if s.RanNgapId != 0 {
		d.RanNgapId = s.RanNgapId
	}

	// GnbId
	if s.GnbId != "" {
		d.GnbId = s.GnbId
	}

	// TacId
	if s.TacId != "" {
		d.TacId = s.TacId
	}

	//	AmfIp
	if s.AmfIp != "" {
		d.AmfIp = s.AmfIp
	}

	// always overwrite subscriber state
	d.AmfSubState = s.AmfSubState
}
