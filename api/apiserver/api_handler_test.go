// Copyright (c) 2026 Intel Corporation
// SPDX-License-Identifier: Apache-2.0

package apiserver

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestWriteJSONResponseSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(recorder)

	if ok := writeJSONResponse(context, gin.H{"value": "ok"}); !ok {
		t.Fatal("expected writeJSONResponse to succeed")
	}
	if recorder.Code != http.StatusOK {
		t.Fatalf("unexpected status: got %d want %d", recorder.Code, http.StatusOK)
	}
	if body := strings.TrimSpace(recorder.Body.String()); body != "{\"value\":\"ok\"}" {
		t.Fatalf("unexpected body: got %q", body)
	}
}

func TestWriteJSONResponseSetBodyFailure(t *testing.T) {
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(recorder)

	if ok := writeJSONResponse(context, map[string]any{"bad": make(chan int)}); ok {
		t.Fatal("expected writeJSONResponse to fail")
	}
	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf("unexpected status: got %d want %d", recorder.Code, http.StatusInternalServerError)
	}
}
