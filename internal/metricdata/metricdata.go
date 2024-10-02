// SPDX-FileCopyrightText: 2022-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package metricdata

import (
	"sync"

	"github.com/omec-project/util/metricinfo"
)

var metricData MetricData

type MetricData struct {
	Subscribers  map[string]*metricinfo.CoreSubscriber
	SubLock      sync.RWMutex
	NfStatusLock sync.RWMutex
	NfStatus     map[string]*metricinfo.CNfStatus
	SmfSvcStats  nfServiceStats
	AmfSvcStats  nfServiceStats
}

func init() {
	metricData = MetricData{
		Subscribers: make(map[string]*metricinfo.CoreSubscriber),
		NfStatus:    make(map[string]*metricinfo.CNfStatus),
		SmfSvcStats: nfServiceStats{svcStats: make(map[string]map[string]uint64)},
		AmfSvcStats: nfServiceStats{svcStats: make(map[string]map[string]uint64)},
	}
}
