// SPDX-FileCopyrightText: 2022-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package metricdata

import (
	"fmt"
	"sync"

	"github.com/omec-project/metricfunc/internal/promclient"
	"github.com/omec-project/metricfunc/logger"
	"github.com/omec-project/util/metricinfo"
)

type nfServiceStats struct {
	svcStatLock sync.RWMutex
	svcStats    map[string]map[string]uint64 // Nf IP is key
}

func HandleServiceEvent(msgType *metricinfo.CoreMsgType, sourceNf metricinfo.NfType) {
	switch sourceNf {
	case metricinfo.NfTypeSmf:
		handleSmfServiceEvent(msgType)
	case metricinfo.NfTypeAmf:
		handleAmfServiceEvent(msgType)
	default:
		logger.CacheLog.Errorf("unknown msg source [%v] ", sourceNf)
	}
}

func handleSmfServiceEvent(msgType *metricinfo.CoreMsgType) {
	metricData.SmfSvcStats.svcStatLock.Lock()
	defer metricData.SmfSvcStats.svcStatLock.Unlock()

	// Check if bucket already exist
	if stats, ok := metricData.SmfSvcStats.svcStats[msgType.SourceNfId]; ok {
		stats[msgType.MsgType] = stats[msgType.MsgType] + 1
		promclient.IncrementSmfSvcStats(msgType.SourceNfId, msgType.MsgType)
		return
	}

	stat := make(map[string]uint64)
	stat[msgType.MsgType] = 1
	metricData.SmfSvcStats.svcStats[msgType.SourceNfId] = stat
	promclient.IncrementSmfSvcStats(msgType.SourceNfId, msgType.MsgType)

	logger.CacheLog.Debugf("smf svc metric data content : %v ", metricData.SmfSvcStats.svcStats)
}

func handleAmfServiceEvent(msgType *metricinfo.CoreMsgType) {
	metricData.AmfSvcStats.svcStatLock.Lock()
	defer metricData.AmfSvcStats.svcStatLock.Unlock()

	// Check if bucket already exist
	if stats, ok := metricData.AmfSvcStats.svcStats[msgType.SourceNfId]; ok {
		stats[msgType.MsgType] = stats[msgType.MsgType] + 1
		promclient.IncrementAmfSvcStats(msgType.SourceNfId, msgType.MsgType)
		return
	}

	stat := make(map[string]uint64)
	stat[msgType.MsgType] = 1
	metricData.AmfSvcStats.svcStats[msgType.SourceNfId] = stat
	promclient.IncrementAmfSvcStats(msgType.SourceNfId, msgType.MsgType)

	logger.CacheLog.Debugf("amf svc metric data content : %v ", metricData.AmfSvcStats.svcStats)
}

func GetNfServiceStatsDetail(nfType string) (map[string](map[string]uint64), error) {
	switch nfType {
	case "smf":
		return GetSmfServiceStatDetail()
	case "amf":
		return GetAmfServiceStatDetail()
	default:
		return nil, fmt.Errorf("no statistics available for nf type [%v] ", nfType)
	}
}

func GetSmfServiceStatDetail() (map[string](map[string]uint64), error) {
	metricData.SmfSvcStats.svcStatLock.RLock()
	defer metricData.SmfSvcStats.svcStatLock.RUnlock()

	return metricData.SmfSvcStats.svcStats, nil
}

func GetAmfServiceStatDetail() (map[string](map[string]uint64), error) {
	metricData.SmfSvcStats.svcStatLock.RLock()
	defer metricData.SmfSvcStats.svcStatLock.RUnlock()

	return metricData.AmfSvcStats.svcStats, nil
}
