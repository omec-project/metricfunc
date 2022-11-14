// SPDX-FileCopyrightText: 2022-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package reader

import (
	"context"
	"encoding/json"
	"log"

	"github.com/omec-project/metricfunc/config"
	"github.com/omec-project/metricfunc/internal/metricdata"
	"github.com/omec-project/metricfunc/pkg/metricinfo"
	"github.com/segmentio/kafka-go"
)

func StartKafkaReader(cfg *config.Configuration) {

	//Start Kafka Event Reader
	for _, nfStream := range cfg.NfStreams {
		r := kafka.NewReader(kafka.ReaderConfig{
			Brokers: nfStream.Urls,
			Topic:   nfStream.Topic.TopicName,
			GroupID: nfStream.Topic.TopicGroups,
		})

		go reader(r)
	}
}

func getSourceNfType(r *kafka.Reader) metricinfo.NfType {

	topic := r.Config().Topic
	switch string(topic[len(topic)-3:]) {
	case "smf":
		return metricinfo.NfTypeSmf
	case "amf":
		return metricinfo.NfTypeAmf
	default:
		log.Default().Fatalf("invalid topic name [%v] ", topic)
		return metricinfo.NfTypeEnd
	}
}

func reader(r *kafka.Reader) {
	log.Printf("kafka reader for topic [%v] initialised ", r.Config().Topic)
	for {
		// the `ReadMessage` function blocks until we receive the next event
		ctxt := context.Background()
		msg, err := r.ReadMessage(ctxt)
		if err != nil {
			log.Printf("Error reading off kafka bus err:%v", err)
		}
		log.Printf("stream [%v] message %s ", r.Config().Topic, string(msg.Value))

		var metricEvent metricinfo.MetricEvent
		//Unmarshal the msg
		if err := json.Unmarshal(msg.Value, &metricEvent); err != nil {
			log.Fatalf("unmarshal smf event error %v ", err.Error())
		}

		sourceNf := getSourceNfType(r)

		switch metricEvent.EventType {
		case metricinfo.CSubscriberEvt:
			metricdata.HandleSubscriberEvent(&metricEvent.SubscriberData, sourceNf)
		case metricinfo.CMsgTypeEvt:
			metricdata.HandleServiceEvent(&metricEvent.MsgType, sourceNf)
		case metricinfo.CNfStatusEvt:
			metricdata.HandleNfStatusEvent(&metricEvent.NfStatusData)
		default:
			log.Fatalf("unknown event type [%v] ", metricEvent.EventType)
		}
	}
}
