// SPDX-FileCopyrightText: 2022-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package apiserver

import (
	"fmt"

	"github.com/omec-project/http2_util"
	"github.com/omec-project/logger_util"
	"github.com/omec-project/metricfunc/config"
	"github.com/omec-project/metricfunc/logger"
)

func init() {
}

func StartApiServer(cfg *config.ServerAddr) {
	router := logger_util.NewGinWithLogrus(logger.GinLog)
	AddService(router)
	HTTPAddr := fmt.Sprintf(":%d", cfg.Port)
	logger.ApiSrvLog.Debugf("api server initialised on address [%v] port [%v] ", cfg.Addr, cfg.Port)
	server, err := http2_util.NewServer(HTTPAddr, "", router)
	if err != nil {
		logger.ApiSrvLog.Errorf("api server initialise error [%v] ", err.Error())
		return
	}

	err = server.ListenAndServe()
	if err != nil {
		logger.ApiSrvLog.Errorf("api server listen error [%v] ", err.Error())
		return
	}
}
