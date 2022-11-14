// SPDX-FileCopyrightText: 2022-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package metricdata

import (
	"fmt"
	"sync"

	"github.com/omec-project/metricfunc/logger"
	"github.com/omec-project/metricfunc/pkg/metricinfo"
)

type nfServiceStats struct {
	svcStatLock sync.RWMutex
	svcStats    map[string]map[string]uint64 //Nf IP is key
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

	//Check if bucket already exist
	if stats, ok := metricData.SmfSvcStats.svcStats[msgType.SourceNfIp]; ok {
		stats[msgType.MsgType] = stats[msgType.MsgType] + 1
		return
	}

	stat := make(map[string]uint64)
	stat[msgType.MsgType] = 1
	metricData.SmfSvcStats.svcStats[msgType.SourceNfIp] = stat

	fmt.Printf("smf svc metric data content : %v ", metricData.SmfSvcStats.svcStats)
}

func handleAmfServiceEvent(msgType *metricinfo.CoreMsgType) {

	metricData.AmfSvcStats.svcStatLock.Lock()
	defer metricData.AmfSvcStats.svcStatLock.Unlock()

	//Check if bucket already exist
	if stats, ok := metricData.AmfSvcStats.svcStats[msgType.SourceNfIp]; ok {
		stats[msgType.MsgType] = stats[msgType.MsgType] + 1
		return
	}

	stat := make(map[string]uint64)
	stat[msgType.MsgType] = 1
	metricData.AmfSvcStats.svcStats[msgType.SourceNfIp] = stat

	fmt.Printf("amf svc metric data content : %v ", metricData.AmfSvcStats.svcStats)
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
	//fillTestSvcStats()

	metricData.SmfSvcStats.svcStatLock.RLock()
	defer metricData.SmfSvcStats.svcStatLock.RUnlock()

	return metricData.SmfSvcStats.svcStats, nil
}

func GetAmfServiceStatDetail() (map[string](map[string]uint64), error) {

	metricData.SmfSvcStats.svcStatLock.RLock()
	defer metricData.SmfSvcStats.svcStatLock.RUnlock()

	return metricData.AmfSvcStats.svcStats, nil
}

/*
func fillTestSvcStats() {
	cm1 := metricinfo.CoreMsgType{MsgType: "pdu_sess_create_req", SourceNfIp: "smf-ip: 1.1.1.1"}
	IncrementSmfServiceStat(&cm1)
	cm12 := metricinfo.CoreMsgType{MsgType: "pdu_sess_update_req", SourceNfIp: "smf-ip: 1.1.1.1"}
	IncrementSmfServiceStat(&cm12)
	cm2 := metricinfo.CoreMsgType{MsgType: "pdu_sess_update_req", SourceNfIp: "smf-ip: 2.2.2.2"}
	IncrementSmfServiceStat(&cm2)
	cm3 := metricinfo.CoreMsgType{MsgType: "pdu_sess_delete_req", SourceNfIp: "smf-ip: 3.3.3.3"}
	IncrementSmfServiceStat(&cm3)
	fmt.Printf("metric data content : %v ", metricData.SmfSvcStats.svcStats)
}
*/
