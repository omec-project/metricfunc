// SPDX-FileCopyrightText: 2026 Open Networking Foundation
//
// SPDX-License-Identifier: Apache-2.0

package apiserver

import (
	"context"
	"testing"
	"time"

	"github.com/omec-project/metricfunc/config"
)

func TestStartApiServerReturnsOnCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() {
		defer close(done)
		StartApiServer(ctx, &config.ServerAddr{Port: 0})
	}()

	time.Sleep(100 * time.Millisecond)
	cancel()
	select {
	case <-done:
	case <-time.After(shutdownTimeout + time.Second):
		t.Fatal("StartApiServer did not return after the context was cancelled")
	}
}
