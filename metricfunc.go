// SPDX-FileCopyrightText: 2022-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"

	"github.com/omec-project/metricfunc/api/apiserver"
	"github.com/omec-project/metricfunc/config"
	"github.com/omec-project/metricfunc/controller"
	"github.com/omec-project/metricfunc/internal/promclient"
	"github.com/omec-project/metricfunc/internal/reader"
	"github.com/omec-project/metricfunc/logger"
	"go.uber.org/zap/zapcore"
	"go.yaml.in/yaml/v4"
)

var PodIp string

func init() {
	podIpStr := os.Getenv("POD_IP")
	PodIp = net.ParseIP(podIpStr).To4().String()
}

func main() {
	// Read provided config
	cfgFilePtr := flag.String("cfg", "/opt/config.yaml", "metricfunc config file")
	flag.Parse()
	logger.AppLog.Infof("Metricfunction has started with configuration file [%v]", *cfgFilePtr)

	cfg := config.Config{}
	if content, err := os.ReadFile(*cfgFilePtr); err != nil {
		logger.AppLog.Errorln("readfile failed", err)
		return
	} else {
		if yamlErr := yaml.Unmarshal(content, &cfg); yamlErr != nil {
			logger.AppLog.Errorln("yaml parsing failed", yamlErr)
			return
		}
	}

	cfg.Configuration.ApiServer.Addr = PodIp
	cfg.Configuration.PrometheusServer.Addr = PodIp

	// set log level
	if level, err := zapcore.ParseLevel(cfg.Logger.LogLevel); err == nil {
		logger.AppLog.Infof("setting log level [%v]", cfg.Logger.LogLevel)
		logger.SetLogLevel(level)
	}

	logger.AppLog.Infof("configuration: %+v", cfg.Configuration)

	// SIGTERM (a pod delete or a rollout) and Ctrl-C stop the servers; main then
	// returns and the process exits 0.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	var servers sync.WaitGroup

	// Start Kafka Event Reader
	reader.StartKafkaReader(cfg.Configuration)

	// Start API Server
	servers.Add(1)
	go func() {
		defer servers.Done()
		apiserver.StartApiServer(ctx, &cfg.Configuration.ApiServer)
	}()

	// Start Prometheus client
	servers.Add(1)
	go func() {
		defer servers.Done()
		promclient.StartPrometheusClient(ctx, &cfg.Configuration.PrometheusServer)
	}()

	if cfg.Configuration.ControllerFlag {
		// controller
		rogueIpChan := make(chan controller.RogueIPs, 100)
		err := controller.InitControllerConfig(&cfg)
		if err != nil {
			logger.AppLog.Warnln("failed to initialize controller configuration")
		}

		userAppClient := controller.UserAppService{
			UserAppServiceUrl: "http://" + cfg.Configuration.UserAppApiServer.Addr + ":" +
				strconv.Itoa(cfg.Configuration.UserAppApiServer.Port) + cfg.Configuration.UserAppApiServer.Path,
			PollInterval: cfg.Configuration.UserAppApiServer.PollInterval,
		}

		controller.RogueChannel = rogueIpChan

		go userAppClient.GetRogueIPs(rogueIpChan)
		go controller.RogueIPHandler(rogueIpChan)
	}

	// Go Pprofiling
	debugProfPort := cfg.Configuration.DebugProfile.Port
	if debugProfPort != 0 {
		logger.AppLog.Infof("pprofile exposed on port [%v]", debugProfPort)
		httpAddr := fmt.Sprintf(":%d", debugProfPort)
		go func() {
			err := http.ListenAndServe(httpAddr, nil)
			if err != nil {
				logger.AppLog.Warnf("failed to listen TCP connection on address %v", httpAddr)
			}
		}()
	}

	<-ctx.Done()
	logger.AppLog.Infoln("terminating metricfunc")
	// The Kafka readers have no consumer group to leave; they end with the
	// process.
	servers.Wait()
	logger.AppLog.Infoln("metricfunc terminated")
}
