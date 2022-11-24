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
	Smf_msg_type_pdu_sess_create_req_failure
	Smf_msg_type_pdu_sess_update_req
	Smf_msg_type_pdu_sess_update_rsp_success
	Smf_msg_type_pdu_sess_update_req_failure
	Smf_msg_type_pdu_sess_release_req
	Smf_msg_type_pdu_sess_release_rsp_success
	Smf_msg_type_pdu_sess_release_req_failure
)

func (t SmfMsgType) String() string {
	switch t {
	case Smf_msg_type_pdu_sess_create_req:
		return "pdu_sess_create_req"
	case Smf_msg_type_pdu_sess_create_rsp_success:
		return "pdu_sess_create_rsp_success"
	case Smf_msg_type_pdu_sess_create_req_failure:
		return "pdu_sess_create_req_failure"
	case Smf_msg_type_pdu_sess_update_req:
		return "pdu_sess_update_req"
	case Smf_msg_type_pdu_sess_update_rsp_success:
		return "pdu_sess_update_rsp_success"
	case Smf_msg_type_pdu_sess_update_req_failure:
		return "pdu_sess_update_req_failure"
	case Smf_msg_type_pdu_sess_release_req:
		return "pdu_sess_release_req"
	case Smf_msg_type_pdu_sess_release_rsp_success:
		return "pdu_sess_release_rsp_success"
	case Smf_msg_type_pdu_sess_release_req_failure:
		return "pdu_sess_release_req_failure"
	default:
		return "unknown smf msg type"
	}
}
