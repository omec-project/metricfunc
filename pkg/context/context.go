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
	CUpf
	CRan
	CAlarm
)

func (e CoreEventType) String() string {
	switch e {
	case CSubscriber:
		return "SUBSCRIBER"
	case CUpf:
		return "UPF"
	case CRan:
		return "RAN"
	case CAlarm:
		return "ALARM"
	}
	return "unknown"
}

// Collective Subscriber info sent towards analytic engine
type CoreSubscriber struct {
	version     int    `json:"version,omitempty"`
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

// UPF state updates
type CoreUpf struct {
	version int    `json:"version,omitempty"`
	UpfName string `json:"upfid"`
	UpfAddr string `json:"upfAddr"`
	State   string `json:"state,omitempty"` // init, Connected, disconnected
}

// RAN state updates
type CoreRan struct {
	version int      `json:"version,omitempty"`
	State   string   `json:"state,omitempty"` // Connected, disconnected
	GnbId   string   `json:"gnbid"`           // string since we can have 0 prefix
	TacId   []string `json:"tacid"`           // string since we can have 0 prefix
	GnbAddr string   `json:"GnbAddr"`         // ipaddress:port
}

// Alarms - UPF DOWN, RAN DOWN, NF Down, DB Down.
type CoreAlarm struct {
	version int `json:"version,omitempty"`
}

type CoreEvent struct {
	version    int             `json:"version,omitempty"`
	Type       string          `json:"type"`
	Subscriber *CoreSubscriber `json:"subscriber,omitempty"`
	Upf        *CoreUpf        `json:"upf,omitempty"`
	Ran        *CoreRan        `json:"ran,omitempty"`
	Alarm      *CoreAlarm      `json:"alarm,omitempty"`
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
