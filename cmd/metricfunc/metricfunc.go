// SPDX-FileCopyrightText: 2022-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"

	//"github.com/omec-project/metricfunc/pkg/consumer"
	ctxt "github.com/omec-project/metricfunc/pkg/context"
	"github.com/omec-project/metricfunc/pkg/producer"
)

type Config struct {
	Info          *Info          `yaml:"info"`
	Logger        *Logger        `yaml:"logger"`
	Configuration *Configuration `yaml:"configuration"`
}

type Info struct {
	Version     string `yaml:"version,omitempty"`
	Description string `yaml:"description,omitempty"`
}

type Logger struct {
	LogLevel string `yaml:"LogLevel,omitempty"`
}

type Configuration struct {
	NfStream        *NFStream        `yaml:"nfStream,omitempty"`
	AnalyticsStream *AnalyticsStream `yaml:"analyticsStream,omitempty"`
}

type Groups struct {
	Analytics  string `yaml: "analytics,omitempty"`
	MongoDB    string `yaml: "mongodb,omitempty"`
	RestApis   string `yaml: "restapi,omitempty"`
	Prometheus string `yaml: "prometheus,omitempty"`
}

type Topic struct {
	Name   string  `yaml: "name,omitempty"`
	Groups *Groups `yaml:"group,omitempty"`
}

type NFStream struct {
	Urls  []string `yaml:"urls,omitempty"`
	Topic *Topic   `yaml:topic,omitempty"`
}

type AnalyticsStream struct {
	Enable    bool     `yaml:"enable,omitempty"`
	Urls      []string `yaml:"urls,omitempty"`
	TopicName string   `yaml:topic,omitempty"`
}

var MetricConfig Config

// TODO : Config updates, logging
func main() {
	log.Println("Metricfunction has started")

	if content, err := ioutil.ReadFile("./config/metricscfg.conf"); err != nil {
		log.Println("Readfile failed called ", err)
		return
	} else {
		MetricConfig = Config{}

		if yamlErr := yaml.Unmarshal(content, &MetricConfig); yamlErr != nil {
			log.Println("yaml parsing failed ", yamlErr)
			return
		}
	}
	if MetricConfig.Configuration == nil {
		log.Println("Configuration Parsing Failed ", MetricConfig.Configuration)
		return
	}

	log.Println("Configuration : ", MetricConfig)

	nf := MetricConfig.Configuration.NfStream

	ec := make(chan error)
	ac := make(chan *ctxt.CoreEvent)
	if nf != nil && len(nf.Topic.Groups.Analytics) > 0 {
		g := nf.Topic.Groups.Analytics
		rt := nf.Topic.Name
		urls := MetricConfig.Configuration.NfStream.Urls
		go nfStreamReader(context.Background(), ac, ec, urls, rt, g)
	}

	/*
		if nf != nil && len(nf.Topic.Groups.MongoDB) > 0 {
			g := nf.Topic.Groups.MongoDB
			rt := nf.Topic.Name
			urls := MetricConfig.Configuration.NfStream.Urls
			go nfStreamReader(context.Background(), mc, ec, urls, rt, g)
		}*/

	as := MetricConfig.Configuration.AnalyticsStream
	if as != nil && as.Enable == true {
		eurls := as.Urls
		et := as.TopicName
		go analyticExporter(context.Background(), ac, ec, eurls, et)
	}

	for {
		select {
		case m := <-ec:
			log.Println("Error Seen ", m)
		}
	}
	return
}

func nfStreamReader(ctx context.Context, ac chan *ctxt.CoreEvent, ec chan error, brokerURLs []string, topic string, groupID string) {
	// initialize a new reader with the brokers and topic
	// the groupID identifies the consumer and prevents
	// it from receiving duplicate messages
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokerURLs,
		Topic:   topic,
		GroupID: groupID,
	})

	brokerStr := "Brokers: "
	for i := 0; i < len(brokerURLs); i++ {
		brokerStr += brokerURLs[i]
	}
	log.Printf("nfStreamReader (%s,%s,%s)", brokerStr, topic, groupID)
	for {
		// the `ReadMessage` function blocks until we receive the next event
		msg, err := r.ReadMessage(ctx)
		if err != nil {
			log.Printf("Error reading off kafka bus err:%v", err)
			ec <- err
			continue
		}
		log.Printf("Received message Group %s %s ", groupID, string(msg.Value))
		e, err := ctxt.GetEvent(msg.Value)
		if err != nil {
			log.Println("Error in parsing event ", err)
			ec <- fmt.Errorf("Error in analyticExporter")
			continue
		}

		if e.Type == "SUBSCRIBER" {
			e.Subscriber = ctxt.StoreSubscriber(e.Subscriber)
			ac <- e
		}
	}
}

func analyticExporter(ctx context.Context, ac chan *ctxt.CoreEvent, ec chan error, brokerURLs []string, topic string) {
	writer := producer.GetWriter(brokerURLs[0], topic)
	for {
		select {
		case msg := <-ac:
			log.Printf("Received Msg in Analytic Export goroutine : %v", msg)
			c, _ := msg.GetMessage()
			log.Printf("Msg for Analytic : %v", msg)
			writer.SendMessage(c)
		}
	}
}
