// SPDX-FileCopyrightText: 2022-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/omec-project/metricfunc/config"
	"github.com/omec-project/metricfunc/internal/metricdata"
	"github.com/omec-project/metricfunc/logger"
	"golang.org/x/net/http2"
)

var (
	ControllerConfig config.Config
	client           *http.Client
)

// creating for testing
var RogueChannel chan RogueIPs

type Targets struct {
	EnterpriseId string `yaml:"name,omitempty" json:"name,omitempty"`
}
type RogueIPs struct {
	IpAddresses []string `yaml:"ipaddresses,omitempty" json:"ipaddresses,omitempty"`
}
type UserAppService struct {
	UserAppServiceUrl string   `yaml:"userAppServiceUrl,omitempty" json:"userAppServiceUrl,omitempty"`
	PollInterval      int      `yaml:"pollInterval,omitempty" json:"pollInterval,omitempty"`
	RogueIPs          RogueIPs `yaml:"rogueips,omitempty" json:"rogueips,omitempty"`
}

type RocService struct {
	RocServiceUrl string     `yaml:"rocServiceUrl,omitempty" json:"rocServiceUrl,omitempty"`
	SiteInfo      []SiteInfo `yaml:"site-info,omitempty" json:"site-info,omitempty"`
}

type SimCard struct {
	SimId       string `yaml:"sim-id,omitempty" json:"sim-id,omitempty"`
	Imsi        string `yaml:"imsi,omitempty" json:"imsi,omitempty"`
	DisplayName string `yaml:"display-name,omitempty" json:"display-name,omitempty"`
	Enable      *bool  `yaml:"enable,omitempty" json:"enable,omitempty"`
}
type SiteInfo struct {
	SiteId         string    `yaml:"site-id,omitempty" json:"site-id,omitempty"`
	SimCardDetails []SimCard `yaml:"sim-card,omitempty" json:"sim-card,omitempty"`
}

func InitControllerConfig(CConfig *config.Config) error {
	ControllerConfig = *CConfig
	// Read provided config
	logger.ControllerLog.Infoln("controller configuration")

	// set http client
	if ControllerConfig.Info.HttpVersion == 2 {
		client = &http.Client{
			Transport: &http2.Transport{
				AllowHTTP: true,
				DialTLS: func(network, addr string, _ *tls.Config) (net.Conn, error) {
					return net.Dial(network, addr)
				},
			},
			Timeout: 5 * time.Second,
		}
	} else {
		client = &http.Client{
			Timeout: 5 * time.Second,
		}
	}

	if ControllerConfig.Configuration.UserAppApiServer.PollInterval == 0 {
		ControllerConfig.Configuration.UserAppApiServer.PollInterval = 30
	}

	logger.ControllerLog.Infoln("ons Api Server Endpoint:")
	addr := ControllerConfig.Configuration.UserAppApiServer.Addr
	ControllerConfig.Configuration.UserAppApiServer.Addr = strings.TrimSpace(addr)
	logger.ControllerLog.Infoln("address", ControllerConfig.Configuration.UserAppApiServer.Addr)
	logger.ControllerLog.Infoln("port", ControllerConfig.Configuration.UserAppApiServer.Port)
	logger.ControllerLog.Infoln("pollInterval", ControllerConfig.Configuration.UserAppApiServer.PollInterval)

	logger.ControllerLog.Infoln("roc Endpoint:")
	ControllerConfig.Configuration.RocEndPoint.Addr = strings.TrimSpace(ControllerConfig.Configuration.RocEndPoint.Addr)
	logger.ControllerLog.Infoln("address", ControllerConfig.Configuration.RocEndPoint.Addr)
	logger.ControllerLog.Infoln("port", ControllerConfig.Configuration.RocEndPoint.Port)

	return nil
}

func getNextBackoffInterval(retry, interval uint) uint {
	mFactor := 1.5
	nextInterval := float64(retry*interval) * mFactor

	if nextInterval > 10 {
		return 10
	}

	return uint(nextInterval)
}

func sendHttpReqMsgWithoutRetry(req *http.Request) (*http.Response, error) {
	rsp, err := client.Do(req)
	if err != nil {
		logger.ControllerLog.Errorf("http req send error [%+v]", err)
		return nil, err
	}

	if rsp.StatusCode == http.StatusAccepted ||
		rsp.StatusCode == http.StatusOK || rsp.StatusCode == http.StatusNoContent ||
		rsp.StatusCode == http.StatusCreated {
		logger.ControllerLog.Infoln("successful response from peer:", http.StatusText(rsp.StatusCode))
		return rsp, nil
	} else {
		logger.ControllerLog.Errorf("http rsp error [%v]", http.StatusText(rsp.StatusCode))
		err := rsp.Body.Close()
		if err != nil {
			logger.ControllerLog.Warnf("body close error: %v", err)
		}

		return nil, fmt.Errorf("error response: %v", http.StatusText(rsp.StatusCode))
	}
}

func sendHttpReqMsg(req *http.Request) (*http.Response, error) {
	// Keep sending request to http server until response is success
	var retries uint = 0
	var body []byte
	var err error
	if req.Body != nil {
		body, err = io.ReadAll(req.Body)
		if err != nil {
			logger.ControllerLog.Warnf("error reading body: %v", err)
		}
	}
	for {
		cloneReq := req.Clone(context.Background())
		req.Body = io.NopCloser(bytes.NewReader(body))
		cloneReq.Body = io.NopCloser(bytes.NewReader(body))
		rsp, err := client.Do(cloneReq)
		retries += 1
		if err != nil {
			nextInterval := getNextBackoffInterval(retries, 2)
			logger.ControllerLog.Warnf("http req send error [%+v], retrying after %v sec...", err, nextInterval)
			time.Sleep(time.Second * time.Duration(nextInterval))
			continue
		}

		if rsp.StatusCode == http.StatusAccepted ||
			rsp.StatusCode == http.StatusOK || rsp.StatusCode == http.StatusNoContent ||
			rsp.StatusCode == http.StatusCreated {
			logger.ControllerLog.Infoln("Get config from peer success")
			err := req.Body.Close()
			if err != nil {
				logger.ControllerLog.Warnf("body close error: %v", err)
			}

			return rsp, nil
		} else {
			nextInterval := getNextBackoffInterval(retries, 2)
			logMsg := "http rsp error [%v], retrying after [%v] sec..."
			logger.ControllerLog.Warnf(logMsg, http.StatusText(rsp.StatusCode), nextInterval)
			err := rsp.Body.Close()
			if err != nil {
				logger.ControllerLog.Warnf("body close error: %v", err)
			}

			time.Sleep(time.Second * time.Duration(nextInterval))
		}
	}
}

func validateIPs(ips RogueIPs) (validIps RogueIPs) {
	for _, ip := range ips.IpAddresses {
		if net.ParseIP(ip) == nil {
			logger.ControllerLog.Errorf("userAppApp response received with IP Address: %s - Invalid", ip)
			continue
		}
		validIps.IpAddresses = append(validIps.IpAddresses, ip)
	}
	logger.ControllerLog.Debugf("rogueIPs [%v] received from userAppApp", validIps.IpAddresses)
	return validIps
}

func (userAppClient *UserAppService) GetRogueIPs(rogueIPChannel chan RogueIPs) {
	userAppServerApi := userAppClient.UserAppServiceUrl
	logger.ControllerLog.Infoln("userAppApp Url:", userAppServerApi)
	req, err := http.NewRequest(http.MethodGet, userAppServerApi, nil)
	if err != nil {
		logger.ControllerLog.Errorln("an error occurred", err)
		return
	}

	for {
		rsp, httpErr := sendHttpReqMsg(req)
		if httpErr != nil {
			logger.ControllerLog.Errorf("get message [%v] returned error [%+v]", userAppServerApi, err)
			time.Sleep(10 * time.Second)
			continue
		}

		var rogueIPs RogueIPs
		if rsp != nil {
			if rsp.Body != nil {
				err := json.NewDecoder(rsp.Body).Decode(&rogueIPs)
				if err != nil {
					logger.ControllerLog.Errorln("userAppApp response body decode failed:", err)
				} else {
					logger.ControllerLog.Infoln("received rogueIPs from userAppApp:", rogueIPs)
					ips := validateIPs(rogueIPs)
					if len(ips.IpAddresses) > 0 {
						// writing rogueIPs into channel
						rogueIPChannel <- ips
					}
				}
			} else {
				logger.ControllerLog.Infoln("http response body from userAppApp is empty")
			}
		}

		time.Sleep(time.Duration(userAppClient.PollInterval) * time.Second)
	}
}

func (rocClient *RocService) GetTargets() (names []Targets) {
	rocTargetsApi := rocClient.RocServiceUrl + "/aether-roc-api/targets"
	req, err := http.NewRequest(http.MethodGet, rocTargetsApi, nil)
	if err != nil {
		logger.ControllerLog.Errorf("get targets request error occurred %v", err)
		return
	}
	rsp, httpErr := sendHttpReqMsgWithoutRetry(req)
	if httpErr != nil {
		logger.ControllerLog.Errorf("get message [%v] returned error [%v] ", rocTargetsApi, httpErr.Error())
		return
	}

	if rsp != nil {
		if rsp.Body != nil {
			err := json.NewDecoder(rsp.Body).Decode(&names)
			if err != nil {
				logger.ControllerLog.Errorln("unable to decode Targets:", err)
			} else {
				logger.ControllerLog.Infoln("get targets received from RoC:", names)
			}
		} else {
			logger.ControllerLog.Errorln("get targets http response body is empty")
		}
	} else {
		logger.ControllerLog.Errorln("get targets http response is empty")
	}
	return
}

func (rocClient *RocService) DisableSimcard(targets []Targets, imsi string) {
	for _, target := range targets {
		rocSiteApi := rocClient.RocServiceUrl + "/aether-roc-api/aether/v2.1.x/" + target.EnterpriseId + "/site"
		req, err := http.NewRequest(http.MethodGet, rocSiteApi, nil)
		if err != nil {
			logger.ControllerLog.Errorf("GetSiteInfo request error occurred %v", err)
			return
		}
		rsp, httpErr := sendHttpReqMsgWithoutRetry(req)
		if httpErr != nil {
			logger.ControllerLog.Errorf("GetSiteInfo message [%v] returned error [%v] ", rocSiteApi, httpErr.Error())
			continue
		}
		var siteInfo []SiteInfo
		if rsp != nil {
			if rsp.Body != nil {
				err := json.NewDecoder(rsp.Body).Decode(&siteInfo)
				if err != nil {
					logger.ControllerLog.Errorln("Unable to decode SiteInfo: ", err)
				} else {
					logger.ControllerLog.Infoln("GetSiteInfo received from RoC: ", siteInfo)
				}

				b, err := io.ReadAll(rsp.Body)
				if err != nil {
					logger.ControllerLog.Warnf("error reading body: %v", err)
				}

				logger.ControllerLog.Infof("SimDetails received from RoC: %s", string(b))
			} else {
				logger.ControllerLog.Errorln("GetSiteInfo http response body is empty")
				continue
			}
		} else {
			logger.ControllerLog.Errorln("GetSiteInfo http response is empty")
			continue
		}

		var rocDisableSimCard *SimCard
		for _, siteInfo := range siteInfo {
			for _, simCard := range siteInfo.SimCardDetails {
				imsi = strings.TrimPrefix(imsi, "imsi-")
				if simCard.Imsi == imsi {
					logger.ControllerLog.Infof("SimCard %v details found in site [%v]", imsi, siteInfo.SiteId)
					rocDisableSimCard = &simCard
					break
				}
			}
			if rocDisableSimCard != nil {
				rocDisableImsiApi := rocSiteApi + "/" + siteInfo.SiteId + "/sim-card/" + rocDisableSimCard.SimId
				var val bool
				rocDisableSimCard.Enable = &val
				b, err := json.Marshal(&rocDisableSimCard)
				if err != nil {
					logger.ControllerLog.Warnf("error marshalling imsi %v: %v", rocDisableSimCard.Imsi, err)
				}

				reqMsgBody := bytes.NewBuffer(b)
				logger.ControllerLog.Debugln("rest API to disable imsi:", rocDisableImsiApi)
				logger.ControllerLog.Debugln("post msg body:", reqMsgBody)

				req, err := http.NewRequest(http.MethodPost, rocDisableImsiApi, reqMsgBody)
				if err != nil {
					logger.ControllerLog.Warnf("error with new request: %v", err)
				}

				req.Header.Set("Content-Type", "application/json; charset=utf-8")
				_, httpErr := sendHttpReqMsgWithoutRetry(req)
				if httpErr != nil {
					logger.ControllerLog.Errorf("post message [%v] returned error [%v]", rocDisableImsiApi, httpErr.Error())
				}
				return
			}
		}
	}

	logger.ControllerLog.Warnf("imsi details not found in Targets and SiteInfo: [%v]", imsi)
}

func RogueIPHandler(rogueIPChannel chan RogueIPs) {
	addr := ControllerConfig.Configuration.RocEndPoint.Addr
	port := ControllerConfig.Configuration.RocEndPoint.Port
	rocClient := RocService{
		RocServiceUrl: "http://" + addr + ":" + strconv.Itoa(port),
	}

	for rogueIPs := range rogueIPChannel {
		for _, ipaddr := range rogueIPs.IpAddresses {
			// get IP to imsi mapping from metricfunc
			subscriberInfo, err := metricdata.GetSubscriberImsiFromIpAddr(ipaddr)
			if err != nil {
				logger.ControllerLog.Errorln("subscriber details doesn't exist with imsi", err)
				continue
			}
			logger.ControllerLog.Infof("subscriber Imsi [%v] of the IP: [%v]", subscriberInfo.Imsi, ipaddr)
			// get enterprises or targets from ROC
			targets := rocClient.GetTargets()

			if len(targets) == 0 {
				logger.ControllerLog.Errorln("get targets returns nil")
			} else {
				// get siteinfo from ROC
				rocClient.DisableSimcard(targets, subscriberInfo.Imsi)
			}
		}
	}
}
