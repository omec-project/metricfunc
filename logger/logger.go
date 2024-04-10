// SPDX-FileCopyrightText: 2022-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package logger

import (
	"time"

	formatter "github.com/antonfisher/nested-logrus-formatter"
	"github.com/sirupsen/logrus"
)

var (
	log           *logrus.Logger
	ApiSrvLog     *logrus.Entry
	GinLog        *logrus.Entry
	CacheLog      *logrus.Entry
	PromLog       *logrus.Entry
	AppLog        *logrus.Entry
	ControllerLog *logrus.Entry
)

func init() {
	log = logrus.New()
	log.SetReportCaller(false)

	log.Formatter = &formatter.Formatter{
		TimestampFormat: time.RFC3339,
		TrimMessages:    true,
		NoFieldsSpace:   true,
		HideKeys:        true,
		FieldsOrder:     []string{"component", "category"},
	}

	ApiSrvLog = log.WithFields(logrus.Fields{"component": "MetricFunc", "category": "ApiServer"})
	GinLog = log.WithFields(logrus.Fields{"component": "MetricFunc", "category": "Gin"})
	CacheLog = log.WithFields(logrus.Fields{"component": "MetricFunc", "category": "Cache"})
	PromLog = log.WithFields(logrus.Fields{"component": "MetricFunc", "category": "Prometheus"})
	AppLog = log.WithFields(logrus.Fields{"component": "MetricFunc", "category": "App"})
	ControllerLog = log.WithFields(logrus.Fields{"component": "Controller", "category": "App"})
}

func SetLogLevel(level logrus.Level) {
	log.SetLevel(level)
}

func SetReportCaller(set bool) {
	log.SetReportCaller(set)
}
