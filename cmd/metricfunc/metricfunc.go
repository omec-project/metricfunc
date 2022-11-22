// SPDX-FileCopyrightText: 2022-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
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

var PodIp string

func init() {
	podIpStr := os.Getenv("POD_IP")
	PodIp = net.ParseIP(podIpStr).To4().String()

}

func main() {

	//Read provided config
	cfgFilePtr := flag.String("metrics", "../../config/config.yaml", "is a config file")
	flag.Parse()
	logger.AppLog.Infof("Metricfunction has started with configuration file [%v]", *cfgFilePtr)

	cfg := config.Config{}
	if content, err := ioutil.ReadFile(*cfgFilePtr); err != nil {
		logger.AppLog.Errorf("Readfile failed called ", err)
		return
	} else {

		if yamlErr := yaml.Unmarshal(content, &cfg); yamlErr != nil {
			logger.AppLog.Errorf("yaml parsing failed ", yamlErr)
			return
		}
	}

	logger.AppLog.Infof("Configuration : %v", cfg.Configuration)
	cfg.Configuration.ApiServer.Addr = PodIp
	cfg.Configuration.PrometheusServer.Addr = PodIp

	//set log level
	if level, err := logrus.ParseLevel(cfg.Logger.LogLevel); err == nil {
		logger.AppLog.Infof("setting log level [%v]", cfg.Logger.LogLevel)
		logger.SetLogLevel(level)
	}

	//Start Kafka Event Reader
	reader.StartKafkaReader(cfg.Configuration)

	//Start API Server
	go apiserver.StartApiServer(&cfg.Configuration.ApiServer)

	//Start Prometheus client
	go promclient.StartPrometheusClient(&cfg.Configuration.PrometheusServer)

	//Go Pprofiling
	debugProfPort := cfg.Configuration.DebugProfile.Port
	if debugProfPort != 0 {
		logger.AppLog.Infof("pprofile exposed on port [%v] ", debugProfPort)
		httpAddr := fmt.Sprintf(":%d", debugProfPort)
		go func() {
			http.ListenAndServe(httpAddr, nil)
		}()
	}

	//Start MongoDB
	select {}
}
