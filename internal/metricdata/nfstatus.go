// SPDX-FileCopyrightText: 2022-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package metricdata

import (
	"github.com/omec-project/metricfunc/internal/promclient"
	"github.com/omec-project/util/metricinfo"
)

func GetNfStatusbyNfType(nfType string) []metricinfo.CNfStatus {
	var nfs []metricinfo.CNfStatus
	metricData.NfStatusLock.RLock()
	defer metricData.NfStatusLock.RUnlock()

	for _, nfStatus := range metricData.NfStatus {
		if nfStatus.NfType == metricinfo.NfType(nfType) {
			nfs = append(nfs, *nfStatus)
		}
	}
	return nfs
}

func GetNfStatusAll() []metricinfo.CNfStatus {
	var nfs []metricinfo.CNfStatus
	metricData.NfStatusLock.RLock()
	defer metricData.NfStatusLock.RUnlock()

	for _, nfStatus := range metricData.NfStatus {
		nfs = append(nfs, *nfStatus)
	}
	return nfs
}

func HandleNfStatusEvent(nfStatus *metricinfo.CNfStatus) {
	metricData.NfStatusLock.Lock()
	defer metricData.NfStatusLock.Unlock()

	metricData.NfStatus[nfStatus.NfName] = nfStatus

	if nfStatus.NfStatus == metricinfo.NfStatusConnected {
		promclient.SetNfStatus(nfStatus.NfName, string(nfStatus.NfType), string(nfStatus.NfStatus), 1)
	} else {
		promclient.SetNfStatus(nfStatus.NfName, string(nfStatus.NfType), string(nfStatus.NfStatus), 0)
	}
}
