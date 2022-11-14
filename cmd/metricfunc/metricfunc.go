// SPDX-FileCopyrightText: 2022-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net"
	"os"

	"net/http"
	_ "net/http/pprof"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"

	"github.com/omec-project/metricfunc/api/apiserver"
	"github.com/omec-project/metricfunc/config"
	"github.com/omec-project/metricfunc/internal/promclient"
	"github.com/omec-project/metricfunc/internal/reader"
	"github.com/omec-project/metricfunc/logger"
)

var TopicsCfg []string
var BrokerURLCfg []string
var PodIp string

func init() {
	TopicsCfg = []string{"sdcore-data-source-smf", "sdcore-data-source-amf"}
	BrokerURLCfg = []string{"sd-core-kafka-headless:9092"}

	podIpStr := os.Getenv("POD_IP")
	PodIp = net.ParseIP(podIpStr).To4().String()

}

func main() {

	//Read provided config
	cfgFilePtr := flag.String("metrics", "../../config/config.yaml", "is a config file")
	flag.Parse()
	log.Printf("Metricfunction has started with configuration file [%v]", *cfgFilePtr)

	cfg := config.Config{}
	if content, err := ioutil.ReadFile(*cfgFilePtr); err != nil {
		log.Println("Readfile failed called ", err)
		return
	} else {

		if yamlErr := yaml.Unmarshal(content, &cfg); yamlErr != nil {
			log.Println("yaml parsing failed ", yamlErr)
			return
		}
	}

	log.Printf("Configuration : %v", cfg.Configuration)
	cfg.Configuration.ApiServer.Addr = PodIp
	cfg.Configuration.PrometheusServer.Addr = PodIp

	//set log level
	if level, err := logrus.ParseLevel(cfg.Logger.LogLevel); err == nil {
		log.Printf("setting log level [%v]", cfg.Logger.LogLevel)
		logger.SetLogLevel(level)
	}

	//Start Kafka Event Reader
	reader.StartKafkaReader(cfg.Configuration)

	//Start API Server
	go apiserver.StartApiServer(&cfg.Configuration.ApiServer)

	//Start Prometheus client
	go promclient.StartPrometheusClient(&cfg.Configuration.PrometheusServer)

	//Go Pprofiling
	go func() {
		http.ListenAndServe(":5001", nil)
	}()

	//Start MongoDB
	select {}
}
