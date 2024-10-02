// SPDX-FileCopyrightText: 2022-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package reader

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/omec-project/metricfunc/config"
	"github.com/omec-project/metricfunc/internal/metricdata"
	"github.com/omec-project/metricfunc/logger"
	"github.com/omec-project/util/metricinfo"
	"github.com/segmentio/kafka-go"
)

func StartKafkaReader(cfg *config.Configuration) {
	// Start Kafka Event Reader
	for _, nfStream := range cfg.NfStreams {
		r := kafka.NewReader(kafka.ReaderConfig{
			Brokers: makeUrlsFromUriPort(nfStream.Urls),
			Topic:   nfStream.Topic.TopicName,
			GroupID: nfStream.Topic.TopicGroups,
		})

		go reader(r)
	}
}

// make urls from config uri and port
func makeUrlsFromUriPort(uriPortCfg []config.Urls) []string {
	var urls []string
	for _, uriPort := range uriPortCfg {
		urls = append(urls, fmt.Sprintf("%s:%d", uriPort.Uri, uriPort.Port))
	}

	return urls
}

func getSourceNfType(r *kafka.Reader) metricinfo.NfType {
	topic := r.Config().Topic
	switch topic {
	case "sdcore-data-source-smf":
		return metricinfo.NfTypeSmf
	case "sdcore-data-source-amf":
		return metricinfo.NfTypeAmf
	default:
		logger.AppLog.Fatalf("invalid topic name [%v]", topic)
		return metricinfo.NfTypeEnd
	}
}

func reader(r *kafka.Reader) {
	logger.AppLog.Infof("kafka reader for topic [%v] initialised", r.Config().Topic)
	sourceNf := getSourceNfType(r)
	for {
		// the `ReadMessage` function blocks until we receive the next event
		ctxt := context.Background()
		msg, err := r.ReadMessage(ctxt)
		if err != nil {
			logger.AppLog.Errorf("Error reading off kafka bus err: %v", err)
			time.Sleep(10 * time.Millisecond)
			continue
		}
		logger.AppLog.Debugf("stream [%v] message %s", r.Config().Topic, string(msg.Value))

		var metricEvent metricinfo.MetricEvent
		// Unmarshal the msg
		if err := json.Unmarshal(msg.Value, &metricEvent); err != nil {
			logger.AppLog.Fatalf("unmarshal smf event error %v", err.Error())
		}

		switch metricEvent.EventType {
		case metricinfo.CSubscriberEvt:
			metricdata.HandleSubscriberEvent(&metricEvent.SubscriberData, sourceNf)
		case metricinfo.CMsgTypeEvt:
			metricdata.HandleServiceEvent(&metricEvent.MsgType, sourceNf)
		case metricinfo.CNfStatusEvt:
			metricdata.HandleNfStatusEvent(&metricEvent.NfStatusData)
		default:
			logger.AppLog.Fatalf("unknown event type [%v]", metricEvent.EventType)
		}
	}
}
