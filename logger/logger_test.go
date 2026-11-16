// Copyright (C) 2026 Intel Corporation
// SPDX-License-Identifier: Apache-2.0

package logger

import (
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// ControllerLog uses its own component name, so it is built apart from the other loggers. It must
// still follow SetLogLevel.
func TestSetLogLevel_AppliesToEveryLogger(t *testing.T) {
	t.Cleanup(func() { component.SetLevel(zapcore.InfoLevel) })

	all := map[string]*zap.SugaredLogger{
		"ApiSrvLog": ApiSrvLog, "GinLog": GinLog, "CacheLog": CacheLog,
		"PromLog": PromLog, "AppLog": AppLog, "ControllerLog": ControllerLog,
	}
	for name, l := range all {
		if l.Desugar().Core().Enabled(zap.DebugLevel) {
			t.Errorf("%s has debug enabled by default", name)
		}
	}

	SetLogLevel(zapcore.DebugLevel)
	for name, l := range all {
		if !l.Desugar().Core().Enabled(zap.DebugLevel) {
			t.Errorf("%s did not follow SetLogLevel", name)
		}
	}
}
