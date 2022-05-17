package producer

import (
	"context"
	"log"
	"github.com/segmentio/kafka-go"
	ctxt "github.com/omec-project/metricfunc/pkg/context"
)

type Writer struct {
	kafkaWriter kafka.Writer
}

/*
GetWriter
creates a kafka.Writer and wraps in a Writer stucture
*/
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

/*
SendMessage
constructs a kafkaMessage from message and writes to
the topic the writer is attached to
*/
func (writer Writer) SendMessage(message []byte) error {
	log.Printf("SendMessage(%s)", string(message))
	msg := kafka.Message{Value: message}
	err := writer.kafkaWriter.WriteMessages(context.Background(), msg)
	return err
}


func PublishCoreSubscriberEvent(writer Writer, sub *ctxt.CoreSubscriber) {
		e := ctxt.GetSubscriberEvent(sub)
		c, _ := e.GetMessage()
		writer.SendMessage(c)
}
