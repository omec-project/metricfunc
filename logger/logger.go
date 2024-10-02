// SPDX-FileCopyrightText: 2022-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	log           *zap.Logger
	ApiSrvLog     *zap.SugaredLogger
	GinLog        *zap.SugaredLogger
	CacheLog      *zap.SugaredLogger
	PromLog       *zap.SugaredLogger
	AppLog        *zap.SugaredLogger
	ControllerLog *zap.SugaredLogger
	atomicLevel   zap.AtomicLevel
)

func init() {
	atomicLevel = zap.NewAtomicLevelAt(zap.InfoLevel)
	config := zap.Config{
		Level:            atomicLevel,
		Development:      false,
		Encoding:         "console",
		EncoderConfig:    zap.NewProductionEncoderConfig(),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncoderConfig.LevelKey = "level"
	config.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	config.EncoderConfig.CallerKey = "caller"
	config.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	config.EncoderConfig.MessageKey = "message"
	config.EncoderConfig.StacktraceKey = ""

	var err error
	log, err = config.Build()
	if err != nil {
		panic(err)
	}

	ApiSrvLog = log.Sugar().With("component", "MetricFunc", "category", "ApiServer")
	GinLog = log.Sugar().With("component", "MetricFunc", "category", "Gin")
	CacheLog = log.Sugar().With("component", "MetricFunc", "category", "Cache")
	PromLog = log.Sugar().With("component", "MetricFunc", "category", "Prometheus")
	AppLog = log.Sugar().With("component", "MetricFunc", "category", "App")
	ControllerLog = log.Sugar().With("component", "Controller", "category", "App")
}

func GetLogger() *zap.Logger {
	return log
}

// SetLogLevel: set the log level (panic|fatal|error|warn|info|debug)
func SetLogLevel(level zapcore.Level) {
	AppLog.Infoln("set log level:", level)
	atomicLevel.SetLevel(level)
}
