// SPDX-FileCopyrightText: 2022-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package apiserver

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/omec-project/metricfunc/config"
	"github.com/omec-project/metricfunc/logger"
	"github.com/omec-project/util/http2_util"
	utilLogger "github.com/omec-project/util/logger"
)

func init() {
}

// shutdownTimeout bounds how long a server waits for requests in flight when
// it is shut down.
const shutdownTimeout = 5 * time.Second

// StartApiServer serves the API until ctx is cancelled, then shuts the server
// down.
func StartApiServer(ctx context.Context, cfg *config.ServerAddr) {
	router := utilLogger.NewGinWithZap(logger.GinLog)
	AddService(router)
	HTTPAddr := fmt.Sprintf(":%d", cfg.Port)
	logger.ApiSrvLog.Debugf("api server initialised on address [%v] port [%v] ", cfg.Addr, cfg.Port)
	server, err := http2_util.NewServer(HTTPAddr, "", router)
	if err != nil {
		logger.ApiSrvLog.Errorf("api server initialise error [%v] ", err.Error())
		return
	}

	serveErr := make(chan error, 1)
	go func() { serveErr <- server.ListenAndServe() }()
	select {
	case err := <-serveErr:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.ApiSrvLog.Errorf("api server listen error [%v] ", err.Error())
		}
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			logger.ApiSrvLog.Warnf("api server shutdown: %v", err)
		}
		logger.ApiSrvLog.Infoln("api server stopped")
	}
}
