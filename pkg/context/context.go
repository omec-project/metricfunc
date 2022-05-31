// SPDX-FileCopyrightText: 2022-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package context

import (
	"encoding/json"
	"fmt"
)

type CoreEventType int64

var Subscribers map[string]*CoreSubscriber

func init() {
	Subscribers = make(map[string]*CoreSubscriber)
}

const (
	CSubscriber CoreEventType = iota
	CNetworkFunction
	CAlarm
)

func (e CoreEventType) String() string {
	switch e {
	case CSubscriber:
		return "SUBSCRIBER"
	case CNetworkFunction:
		return "NF"
	case CAlarm:
		return "ALARM"
	}
	return "Unknown"
}

// Collective Subscriber info sent towards analytic engine
type CoreSubscriber struct {
	Version     int    `json:"version,omitempty"`
	Imsi        string `json:"imsi, omitempty"`
	SmfId       string `json:"smfId,omitempty"`
	SmfSubState string `json:"smfSubState,omitempty"` //Connected, Idle, Deleted
	IPAddress   string `json:"ipaddress, omitempty"`
	LSEID       int    `json:"lseid,omitempty"`
	RSEID       int    `json:"rseid,omitempty"`
	UpfName     string `json:"upfid"`
	UpfAddr     string `json:"upfAddr"`
	AmfId       string `json:"amfId,omitempty"`
	Guti        string `json:"guti,omitempty"`
	Tmsi        string `json:"tmsi,omitempty"`
	AmfNgapId   int    `json:"amfngapId,omitempty"`
	RanNgapId   int    `json:"ranngapId,omitempty"`
	AmfSubState string `json:"amfSubState, omitempty"` //RegisteredC, RegisteredI, Deregistered, Deleted
	GnbId       string `json:"gnbid"`
	TacId       string `json:"tacid"`
}

type CoreNetworkFunction struct {
	Version int      `json:"version,omitempty"`
	nfType  string   `json:"nftype,omitempty"` // gNB, UPF, AMF, UPF, SMF, NRF,...
	State   string   `json:"state,omitempty"`  // Unknown, Active, Disconnected
	nfId    string   `json:"nfid,omitempty"`            // string since we can have 0 prefix
	nfName  string   `json:"name,omitempty"`
	TacId   []string `json:"tacid,omitempty"`  // string since we can have 0 prefix
	nfAddr  string   `json:"nfAddr,omitempty"` // ipaddress:port
}

// Alarms - UPF DOWN, RAN DOWN, NF Down, DB Down.
type CoreAlarm struct {
	Version     int    `json:"version,omitempty"`
	Description string `json: "description,omitempty"`
}

type CoreEvent struct {
	Version    int              `json:"version,omitempty"`
	Type       string           `json:"type"`
	Subscriber *CoreSubscriber  `json:"subscriber,omitempty"`
	NetworkFn  *CoreNetworkFunction `json:"networkfunction,omitempty"`
	Alarm      *CoreAlarm       `json:"alarm,omitempty"`
}

func GetSubscriberEvent(sub *CoreSubscriber) CoreEvent {
	s := CSubscriber
	st := fmt.Sprintf("%s", s)
	e := CoreEvent{Type: st, Subscriber: sub}
	return e
}

func StoreSubscriber(sub *CoreSubscriber) *CoreSubscriber {
	s, ok := Subscribers[sub.Imsi]
	if !ok {
		Subscribers[sub.Imsi] = sub
		s = sub
	} else {
		if len(sub.IPAddress) > 0 {
			s.IPAddress = sub.IPAddress
		}
		if len(sub.SmfSubState) > 0 {
			s.SmfSubState = sub.SmfSubState
		}
	}
	return s
}

func GetNetworkFunctionEvent(nf *CoreNetworkFunction) CoreEvent {
	s := CNetworkFunction
	st := fmt.Sprintf("%s", s)
	e := CoreEvent{Type: st, NetworkFn: nf}
	return e
}

func GetAlarmEvent(al *CoreAlarm) CoreEvent {
	s := CAlarm
	st := fmt.Sprintf("%s", s)
	e := CoreEvent{Type: st, Alarm: al}
	return e
}

func (e CoreEvent) GetMessage() ([]byte, error) {
	m, err := json.Marshal(e)
	if err != nil {
		fmt.Println("String Marshalled error ", err)
		return []byte{}, err
	}
	return m, nil
}

func GetEvent(b []byte) (*CoreEvent, error) {
	e := &CoreEvent{}
	err := json.Unmarshal(b, e)
	if err != nil {
		fmt.Println("String Marshalled error ", err)
		return e, err
	}
	return e, nil
}
