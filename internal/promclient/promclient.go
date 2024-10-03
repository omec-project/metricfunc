// SPDX-FileCopyrightText: 2022-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package promclient

import (
	"fmt"
	"net/http"

	"github.com/omec-project/metricfunc/config"
	"github.com/omec-project/metricfunc/logger"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type PromStats struct {
	coreSub     *prometheus.CounterVec
	smfSvcStat  *prometheus.CounterVec
	amfSvcStat  *prometheus.CounterVec
	smfSessions *prometheus.GaugeVec
	nfStatus    *prometheus.GaugeVec
}

var promStats *PromStats

func init() {
	promStats = initPromStats()

	if err := promStats.register(); err != nil {
		logger.PromLog.Panicln("prometheus stats register failed", err.Error())
	}
}

func StartPrometheusClient(cfg *config.ServerAddr) {
	logger.PromLog.Debugf("prometheus server initialised on address [%v] port [%v]", cfg.Addr, cfg.Port)
	HTTPAddr := fmt.Sprintf(":%d", cfg.Port)
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(HTTPAddr, nil)
}

func initPromStats() *PromStats {
	return &PromStats{
		coreSub: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "core_subscriber",
			Help: "core subscriber info",
		}, []string{"imsi", "ip_addr", "state", "smf_ip", "dnn", "slice", "upf"}),

		smfSessions: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "smf_pdu_sessions",
			Help: "Number of SMF PDU sessions currently in the core",
		}, []string{"smf_ip", "slice", "dnn", "upf"}),

		nfStatus: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "nf_status",
			Help: "NF Status up/down",
		}, []string{"Nfname", "nfType"}),

		smfSvcStat: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "smf_svc_stats",
			Help: "smf service stats",
		}, []string{"smfid", "msgtype"}),

		amfSvcStat: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "amf_svc_stats",
			Help: "amf service stats",
		}, []string{"amfid", "msgtype"}),
	}
}

func (ps *PromStats) register() error {
	if err := prometheus.Register(ps.coreSub); err != nil {
		logger.PromLog.Errorf("register core subscriber detail stats failed: %v", err.Error())
		return err
	}

	if err := prometheus.Register(ps.smfSessions); err != nil {
		logger.PromLog.Errorf("register core subscriber count stats failed: %v", err.Error())
		return err
	}

	if err := prometheus.Register(ps.nfStatus); err != nil {
		logger.PromLog.Errorf("register nf status stats failed: %v", err.Error())
		return err
	}

	if err := prometheus.Register(ps.smfSvcStat); err != nil {
		logger.PromLog.Errorf("register smf service stats failed: %v", err.Error())
		return err
	}

	if err := prometheus.Register(ps.amfSvcStat); err != nil {
		logger.PromLog.Errorf("register amf service stats failed: %v", err.Error())
		return err
	}
	return nil
}

// PushCoreSubData increments message level stats
func PushCoreSubData(imsi, ip_addr, state, smf_ip, dnn, slice, upf string) {
	logger.PromLog.Debugf("adding subscriber data [%v, %v, %v, %v, %v, %v, %v]", imsi, ip_addr, state, smf_ip, dnn, slice, upf)
	promStats.coreSub.WithLabelValues(imsi, "", state, smf_ip, dnn, slice, upf).Inc()
}

func DeleteCoreSubData(imsi, ip_addr, state, smf_ip, dnn, slice, upf string) {
	logger.PromLog.Debugf("deleting subscriber data [%v, %v, %v, %v, %v, %v, %v]", imsi, ip_addr, state, smf_ip, dnn, slice, upf)
	promStats.coreSub.DeleteLabelValues(imsi, "", state, smf_ip, dnn, slice, upf)
}

// SetSessStats maintains Session level stats
func SetSmfSessStats(smfIp, slice, dnn, upf string, count uint64) {
	logger.PromLog.Debugf("setting smf session count [%v] with labels [smfIp:%v, slice:%v, dnn:%v, upf:%v]", count, smfIp, slice, dnn, upf)
	promStats.smfSessions.WithLabelValues("", "", "", "").Set(float64(count))
}

func DeleteSmfSessStats(smfIp, slice, dnn, upf string) {
	logger.PromLog.Warnln("deleting smf session count stats")
	promStats.smfSessions.DeleteLabelValues(smfIp, slice, dnn, upf)
}

func SetNfStatus(nfName, nfType, nfStatus string, value uint64) {
	logger.PromLog.Debugf("setting nf [%v], type [%v],  status to [%v]", nfName, nfType, nfStatus)
	promStats.nfStatus.WithLabelValues(nfName, nfType).Set(float64(value))
}

func IncrementSmfSvcStats(smfId, msgType string) {
	logger.PromLog.Debugf("incrementing smf service stats, instance [%v] msgtype [%v]", smfId, msgType)
	promStats.smfSvcStat.WithLabelValues(smfId, msgType).Inc()
}

func IncrementAmfSvcStats(amfId, msgType string) {
	logger.PromLog.Debugf("incrementing amf service stats, instance [%v] msgtype [%v]", amfId, msgType)
	promStats.smfSvcStat.WithLabelValues(amfId, msgType).Inc()
}
