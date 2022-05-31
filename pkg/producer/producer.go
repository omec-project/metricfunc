// SPDX-FileCopyrightText: 2022-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package producer

import (
	"context"
	ctxt "github.com/omec-project/metricfunc/pkg/context"
	"github.com/segmentio/kafka-go"
	"log"
)

type Writer struct {
	kafkaWriter kafka.Writer
}

func GetWriter(kafkaURI string, topic string) Writer {
	log.Printf("GetWriter(%s,%s)", kafkaURI, topic)
	producer := kafka.Writer{
		Addr:     kafka.TCP(kafkaURI),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
	writer := Writer{
		kafkaWriter: producer,
	}
	return writer
}

func (writer Writer) SendMessage(message []byte) error {
	msg := kafka.Message{Value: message}
	err := writer.kafkaWriter.WriteMessages(context.Background(), msg)
	return err
}

func PublishCoreSubscriberEvent(writer Writer, sub *ctxt.CoreSubscriber) {
	e := ctxt.GetSubscriberEvent(sub)
	c, _ := e.GetMessage()
	log.Printf("Publish CoreSubscriber Message(%s)", string(c))
	writer.SendMessage(c)
}

func PublishPeerNFEvent(writer Writer, nf *ctxt.CoreNetworkFunction) {
	e := ctxt.GetNetworkFunctionEvent(nf)
	c, _ := e.GetMessage()
	log.Printf("Publish CoreNetworkFunction Message(%s)", string(c))
	writer.SendMessage(c)
}

func PublishAlarmEvent(writer Writer, al *ctxt.CoreAlarm) {
	e := ctxt.GetAlarmEvent(al)
	c, _ := e.GetMessage()
	log.Printf("Publish CoreAlarm Message(%s)", string(c))
	writer.SendMessage(c)
}
