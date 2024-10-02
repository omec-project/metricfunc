// SPDX-FileCopyrightText: 2022-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package metricdata

import (
	"fmt"
	"sync/atomic"

	"github.com/omec-project/metricfunc/internal/promclient"
	"github.com/omec-project/metricfunc/logger"
	"github.com/omec-project/util/metricinfo"
)

var smContextActive uint64

func incSMContextActive() uint64 {
	atomic.AddUint64(&smContextActive, 1)
	return smContextActive
}

func decSMContextActive() uint64 {
	atomic.AddUint64(&smContextActive, ^uint64(0))
	return smContextActive
}

func HandleSubscriberEvent(subsData *metricinfo.CoreSubscriberData, sourceNf metricinfo.NfType) {
	switch subsData.Operation {
	case metricinfo.SubsOpAdd:
		err := storeSubscriber(&subsData.Subscriber, sourceNf)
		if err != nil {
			logger.CacheLog.Infof("store subsriber %v failed for sourceNF [%v] ", subsData.Subscriber.Imsi, sourceNf)
		}
	case metricinfo.SubsOpMod:
		updateSubscriber(&subsData.Subscriber, sourceNf)
	case metricinfo.SubsOpDel:
		err := deleteSubscriber(&subsData.Subscriber, sourceNf)
		if err != nil {
			logger.CacheLog.Infof("store subsriber %v failed for sourceNF [%v] ", subsData.Subscriber.Imsi, sourceNf)
		}
	default:
		logger.CacheLog.Errorf("unknown smf subsriber operation [%v] ", subsData.Operation)
	}
}

func storeSubscriber(sub *metricinfo.CoreSubscriber, sourceNf metricinfo.NfType) error {
	metricData.SubLock.Lock()

	_, ok := metricData.Subscribers[sub.Imsi]
	if !ok {
		metricData.Subscribers[sub.Imsi] = sub

		promclient.SetSmfSessStats(sub.SmfIp, sub.Slice, sub.Dnn, sub.UpfName, incSMContextActive())
		logger.CacheLog.Debugf("storing subscriber with imsi [%s] ", sub.Imsi)
		pushPrometheusCoreSubData(sub)
		metricData.SubLock.Unlock()
	} else {
		metricData.SubLock.Unlock()
		updateSubscriber(sub, sourceNf)
	}

	return nil
}

func updateSubscriber(sub *metricinfo.CoreSubscriber, sourceNf metricinfo.NfType) {
	metricData.SubLock.Lock()
	defer metricData.SubLock.Unlock()
	if s, ok := metricData.Subscribers[sub.Imsi]; ok {
		deletePrometheusCoreSubData(s)

		if sourceNf == metricinfo.NfTypeSmf {
			// SMF specific fields
			fillSmfSubsriberData(sub, s)
		} else if sourceNf == metricinfo.NfTypeAmf {
			// AMF specific fields
			fillAmfSubsriberData(sub, s)
		}
		pushPrometheusCoreSubData(s)
	}
}

func deleteSubscriber(sub *metricinfo.CoreSubscriber, sourceNf metricinfo.NfType) error {
	metricData.SubLock.Lock()
	defer metricData.SubLock.Unlock()
	imsi := sub.Imsi
	s, ok := metricData.Subscribers[imsi]
	if !ok {
		return fmt.Errorf("subscriber with imsi [%s] already deleted ", imsi)
	}

	promclient.SetSmfSessStats(s.SmfIp, s.Slice, s.Dnn, s.UpfName, decSMContextActive())
	deletePrometheusCoreSubData(s)
	s.SmfSubState = sub.SmfSubState
	s.AmfSubState = sub.AmfSubState

	// register disconnect state
	pushPrometheusCoreSubData(s)
	delete(metricData.Subscribers, imsi)

	// register subscriber delete
	deletePrometheusCoreSubData(s)

	logger.CacheLog.Debugf("deleting subscriber with imsi [%s] ", imsi)

	return nil
}

func GetSubscriber(key string) (*metricinfo.CoreSubscriber, error) {
	metricData.SubLock.RLock()
	defer metricData.SubLock.RUnlock()
	if sub, ok := metricData.Subscribers[key]; ok {
		return sub, nil
	}
	return nil, fmt.Errorf("subscriber with key [%v] not found ", key)
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
	return nil, fmt.Errorf("subscriber with ip-addr [%v] not found ", ipaddr)
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
