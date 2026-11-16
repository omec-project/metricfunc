// Copyright (C) 2022-Present Intel Corporation
// SPDX-License-Identifier: Apache-2.0

package logger

import (
	utilLogger "github.com/omec-project/util/logger"
	"go.uber.org/zap/zapcore"
)

var component = utilLogger.NewComponent("MetricFunc")

var (
	ApiSrvLog = component.Category("ApiServer")
	GinLog    = component.Category("Gin")
	CacheLog  = component.Category("Cache")
	PromLog   = component.Category("Prometheus")
	AppLog    = component.Category("App")

	// ControllerLog logs under another component name but shares the logger, and so the level, of the rest.
	ControllerLog = component.Zap().Sugar().With("component", "Controller", "category", "App")
)

// SetLogLevel: set the log level (panic|fatal|error|warn|info|debug)
func SetLogLevel(level zapcore.Level) {
	AppLog.Infoln("set log level:", level)
	component.SetLevel(level)
}
