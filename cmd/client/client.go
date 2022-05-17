package main

import (
	"github.com/omec-project/metricfunc/pkg/context"
	"github.com/omec-project/metricfunc/pkg/producer"
	"log"
	"strconv"
	"time"
)

// Configuration needs :
//  Kafka endpoint
//  topic
func main() {
	log.Println("Metricgen has started")
	writer := producer.GetWriter("sd-core-kafka-headless:9092", "sdcore-nf-data-source")

	count := 123456
	for {
		log.Println("Client Iteration ", count)
		count = count + 1
		imsi := strconv.Itoa(count)
		sub := &context.CoreSubscriber{Imsi: imsi, SmfSubState: "Connected", IPAddress: "1.1.1.1"}
		producer.PublishCoreSubscriberEvent(writer, sub);
		time.Sleep(10 * time.Second)
	}
	return
}
