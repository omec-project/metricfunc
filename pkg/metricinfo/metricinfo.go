// SPDX-FileCopyrightText: 2022-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package metricinfo

type CoreSubscriber struct {
	Version     int    `json:"version,omitempty"`
	Imsi        string `json:"imsi,omitempty"` //key
	SmfId       string `json:"smfId,omitempty"`
	SmfIp       string `json:"smfIp,omitempty"`
	SmfSubState string `json:"smfSubState,omitempty"` //Connected, Idle, DisConnected
	IPAddress   string `json:"ipaddress,omitempty"`
	Dnn         string `json:"dnn,omitempty"`
	Slice       string `json:"slice,omitempty"`
	LSEID       int    `json:"lseid,omitempty"`
	RSEID       int    `json:"rseid,omitempty"`
	UpfName     string `json:"upfid,omitempty"`
	UpfAddr     string `json:"upfAddr,omitempty"`
	AmfId       string `json:"amfId,omitempty"`
	Guti        string `json:"guti,omitempty"`
	Tmsi        int32  `json:"tmsi,omitempty"`
	AmfNgapId   int64  `json:"amfngapId,omitempty"`
	RanNgapId   int64  `json:"ranngapId,omitempty"`
	AmfSubState string `json:"amfSubState,omitempty"` //RegisteredC, RegisteredI, DeRegistered, Deleted
	GnbId       string `json:"gnbid,omitempty"`
	TacId       string `json:"tacid,omitempty"`
	AmfIp       string `json:"amfIp,omitempty"`
	UeState     string `json:"ueState,omitempty"`
}

type CoreMsgType struct {
	MsgType    string `json:"msgType,omitempty"`
	SourceNfId string `json:"sourceNfId,omitempty"`
}

type CoreEventType int64

const (
	CSubscriberEvt CoreEventType = iota
	CMsgTypeEvt
	CNfStatusEvt
)

func (e CoreEventType) String() string {
	switch e {
	case CSubscriberEvt:
		return "SubscriberEvt"
	case CMsgTypeEvt:
		return "MsgTypeEvt"
	case CNfStatusEvt:
		return "CNfStatusEvt"
	}
	return "Unknown"
}

type NfStatusType string

const (
	NfStatusConnected    NfStatusType = "Connected"
	NfStatusDisconnected NfStatusType = "Disconnected"
)

type NfType string

const (
	NfTypeSmf NfType = "SMF"
	NfTypeAmf NfType = "AMF"
	NfTypeUPF NfType = "UPF"
	NfTypeGnb NfType = "GNB"
	NfTypeEnd NfType = "Invalid"
)

type CNfStatus struct {
	NfType   NfType       `json:"nfType,omitempty"`
	NfStatus NfStatusType `json:"nfStatus,omitempty"`
	NfName   string       `json:"nfName,omitempty"`
}

type SubscriberOp uint

const (
	SubsOpAdd SubscriberOp = iota + 1
	SubsOpMod
	SubsOpDel
)

type CoreSubscriberData struct {
	Subscriber CoreSubscriber `json:"subscriber,omitempty"`
	Operation  SubscriberOp   `json:"subsOp,omitempty"`
}

//Sent by NFs(Producers) and received by Metric Function
type MetricEvent struct {
	EventType      CoreEventType      `json:"eventType,omitempty"`
	SubscriberData CoreSubscriberData `json:"subscriberData,omitempty"`
	MsgType        CoreMsgType        `json:"coreMsgType,omitempty"`
	NfStatusData   CNfStatus          `json:"nfStatusData"`
}

type SmfMsgType uint64

const (
	Smf_msg_type_invalid SmfMsgType = iota
	Smf_msg_type_pdu_sess_create_req
	Smf_msg_type_pdu_sess_create_rsp_success
	Smf_msg_type_pdu_sess_create_rsp_failure
	Smf_msg_type_pdu_sess_modify_req
	Smf_msg_type_pdu_sess_modify_rsp_success
	Smf_msg_type_pdu_sess_modify_rsp_failure
	Smf_msg_type_pdu_sess_release_req
	Smf_msg_type_pdu_sess_release_rsp_success
	Smf_msg_type_pdu_sess_release_rsp_failure
	Smf_msg_type_n1n2_transfer_req
	Smf_msg_type_n1n2_transfer_rsp_success
	Smf_msg_type_n1n2_transfer_rsp_failure
	Smf_msg_type_smpolicy_create_req
	Smf_msg_type_smpolicy_create_rsp_success
	Smf_msg_type_smpolicy_create_rsp_failure
	Smf_msg_type_pfcp_sess_estab_req
	Smf_msg_type_pfcp_sess_estab_rsp_success
	Smf_msg_type_pfcp_sess_estab_rsp_failure
	Smf_msg_type_pfcp_sess_modify_req
	Smf_msg_type_pfcp_sess_modify_rsp_success
	Smf_msg_type_pfcp_sess_modify_rsp_failure
	Smf_msg_type_pfcp_sess_release_req
	Smf_msg_type_pfcp_sess_release_rsp_success
	Smf_msg_type_pfcp_sess_release_rsp_failure
	Smf_msg_type_pfcp_association_req
	Smf_msg_type_pfcp_association_rsp_success
	Smf_msg_type_pfcp_association_rsp_failure
	Smf_msg_type_pfcp_heartbeat_req
	Smf_msg_type_pfcp_heartbeat_rsp_success
	Smf_msg_type_pfcp_heartbeat_rsp_failure
	Smf_msg_type_udm_get_smdata_req
	Smf_msg_type_udm_get_smdata_rsp_success
	Smf_msg_type_udm_get_smdata_rsp_failure
	Smf_msg_type_nrf_discovery_amf_req
	Smf_msg_type_nrf_discovery_amf_rsp_success
	Smf_msg_type_nrf_discovery_amf_rsp_failure
	Smf_msg_type_nrf_discovery_pcf_req
	Smf_msg_type_nrf_discovery_pcf_rsp_success
	Smf_msg_type_nrf_discovery_pcf_rsp_failure
	Smf_msg_type_nrf_discovery_udm_req
	Smf_msg_type_nrf_discovery_udm_rsp_success
	Smf_msg_type_nrf_discovery_udm_rsp_failure
	Smf_msg_type_nrf_register_smf_req
	Smf_msg_type_nrf_register_smf_rsp_success
	Smf_msg_type_nrf_register_smf_rsp_failure
	Smf_msg_type_nrf_deregister_smf_req
	Smf_msg_type_nrf_deregister_smf_rsp_success
	Smf_msg_type_nrf_deregister_smf_rsp_failure
)

func (t SmfMsgType) String() string {
	switch t {
	case Smf_msg_type_pdu_sess_create_req:
		return "pdu_sess_create_req"
	case Smf_msg_type_pdu_sess_create_rsp_success:
		return "pdu_sess_create_rsp_success"
	case Smf_msg_type_pdu_sess_create_rsp_failure:
		return "pdu_sess_create_rsp_failure"
	case Smf_msg_type_pdu_sess_modify_req:
		return "pdu_sess_modify_req"
	case Smf_msg_type_pdu_sess_modify_rsp_success:
		return "pdu_sess_modify_rsp_success"
	case Smf_msg_type_pdu_sess_modify_rsp_failure:
		return "pdu_sess_modify_rsp_failure"
	case Smf_msg_type_pdu_sess_release_req:
		return "pdu_sess_release_req"
	case Smf_msg_type_pdu_sess_release_rsp_success:
		return "pdu_sess_release_rsp_success"
	case Smf_msg_type_pdu_sess_release_rsp_failure:
		return "pdu_sess_release_rsp_failure"
	case Smf_msg_type_n1n2_transfer_req:
		return "n1n2_transfer_req"
	case Smf_msg_type_n1n2_transfer_rsp_success:
		return "n1n2_transfer_rsp_success"
	case Smf_msg_type_n1n2_transfer_rsp_failure:
		return "n1n2_transfer_rsp_failure"
	case Smf_msg_type_smpolicy_create_req:
		return "smpolicy_create_req"
	case Smf_msg_type_smpolicy_create_rsp_success:
		return "smpolicy_create_rsp_success"
	case Smf_msg_type_smpolicy_create_rsp_failure:
		return "smpolicy_create_rsp_failure"
	case Smf_msg_type_pfcp_sess_estab_req:
		return "pfcp_sess_estab_req"
	case Smf_msg_type_pfcp_sess_estab_rsp_success:
		return "pfcp_sess_estab_rsp_success"
	case Smf_msg_type_pfcp_sess_estab_rsp_failure:
		return "pfcp_sess_estab_rsp_failure"
	case Smf_msg_type_pfcp_sess_modify_req:
		return "pfcp_sess_modify_req"
	case Smf_msg_type_pfcp_sess_modify_rsp_success:
		return "pfcp_sess_modify_rsp_success"
	case Smf_msg_type_pfcp_sess_modify_rsp_failure:
		return "pfcp_sess_modify_rsp_failure"
	case Smf_msg_type_pfcp_sess_release_req:
		return "pfcp_sess_release_req"
	case Smf_msg_type_pfcp_sess_release_rsp_success:
		return "pfcp_sess_release_rsp_success"
	case Smf_msg_type_pfcp_sess_release_rsp_failure:
		return "pfcp_sess_release_rsp_failure"
	case Smf_msg_type_pfcp_association_req:
		return "pfcp_association_req"
	case Smf_msg_type_pfcp_association_rsp_success:
		return "pfcp_association_rsp_success"
	case Smf_msg_type_pfcp_association_rsp_failure:
		return "pfcp_association_rsp_failure"
	case Smf_msg_type_pfcp_heartbeat_req:
		return "pfcp_heartbeat_req"
	case Smf_msg_type_pfcp_heartbeat_rsp_success:
		return "pfcp_heartbeat_rsp_success"
	case Smf_msg_type_pfcp_heartbeat_rsp_failure:
		return "pfcp_heartbeat_rsp_failure"
	case Smf_msg_type_udm_get_smdata_req:
		return "udm_get_smdata_req"
	case Smf_msg_type_udm_get_smdata_rsp_success:
		return "udm_get_smdata_rsp_success"
	case Smf_msg_type_udm_get_smdata_rsp_failure:
		return "udm_get_smdata_rsp_failure"
	case Smf_msg_type_nrf_discovery_amf_req:
		return "nrf_discovery_amf_req"
	case Smf_msg_type_nrf_discovery_amf_rsp_success:
		return "nrf_discovery_amf_rsp_success"
	case Smf_msg_type_nrf_discovery_amf_rsp_failure:
		return "nrf_discovery_amf_rsp_failure"
	case Smf_msg_type_nrf_discovery_pcf_req:
		return "nrf_discovery_pcf_req"
	case Smf_msg_type_nrf_discovery_pcf_rsp_success:
		return "nrf_discovery_pcf_rsp_success"
	case Smf_msg_type_nrf_discovery_pcf_rsp_failure:
		return "nrf_discovery_pcf_rsp_failure"
	case Smf_msg_type_nrf_discovery_udm_req:
		return "nrf_discovery_udm_req"
	case Smf_msg_type_nrf_discovery_udm_rsp_success:
		return "nrf_discovery_udm_rsp_success"
	case Smf_msg_type_nrf_discovery_udm_rsp_failure:
		return "nrf_discovery_udm_rsp_failure"
	case Smf_msg_type_nrf_register_smf_req:
		return "nrf_register_smf_req"
	case Smf_msg_type_nrf_register_smf_rsp_success:
		return "nrf_register_smf_rsp_success"
	case Smf_msg_type_nrf_register_smf_rsp_failure:
		return "nrf_register_smf_rsp_failure"
	case Smf_msg_type_nrf_deregister_smf_req:
		return "nrf_deregister_smf_req"
	case Smf_msg_type_nrf_deregister_smf_rsp_success:
		return "nrf_deregister_smf_rsp_success"
	case Smf_msg_type_nrf_deregister_smf_rsp_failure:
		return "nrf_deregister_smf_rsp_failure"
	default:
		return "unknown smf msg type"
	}
}
