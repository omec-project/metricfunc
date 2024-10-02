// SPDX-FileCopyrightText: 2022-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package metricdata

import (
	"github.com/omec-project/metricfunc/logger"
	"github.com/omec-project/util/metricinfo"
)

func TestfillTestSmfSvcStats() {
	cm1 := metricinfo.CoreMsgType{MsgType: "pdu_sess_create_req", SourceNfId: "smf-ip: 1.1.1.1"}
	handleSmfServiceEvent(&cm1)
	cm12 := metricinfo.CoreMsgType{MsgType: "pdu_sess_update_req", SourceNfId: "smf-ip: 1.1.1.1"}
	handleSmfServiceEvent(&cm12)
	cm2 := metricinfo.CoreMsgType{MsgType: "pdu_sess_update_req", SourceNfId: "smf-ip: 2.2.2.2"}
	handleSmfServiceEvent(&cm2)
	cm3 := metricinfo.CoreMsgType{MsgType: "pdu_sess_delete_req", SourceNfId: "smf-ip: 3.3.3.3"}
	handleSmfServiceEvent(&cm3)
	logger.CacheLog.Debugf("metric data content : %v ", metricData.SmfSvcStats.svcStats)
}
