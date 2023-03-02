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
	"io/ioutil"
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

var ControllerConfig config.Config
var client *http.Client

//creating for testing
var RogueChannel chan RogueIPs

type Targets struct {
	EnterpriseId string `yaml:"name,omitempty" json:"name,omitempty"`
}
type RogueIPs struct {
	IpAddresses []string `yaml:"ipaddresses,omitempty" json:"ipaddresses,omitempty"`
}
type OnosService struct {
	OnosServiceUrl string   `yaml:"onosServiceUrl,omitempty" json:"onosServiceUrl,omitempty"`
	PollInterval   int      `yaml:"pollInterval,omitempty" json:"pollInterval,omitempty"`
	RogueIPs       RogueIPs `yaml:"rogueips,omitempty" json:"rogueips,omitempty"`
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
	//Read provided config
	fmt.Printf("Controller configuration")

	//set http client
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

	if ControllerConfig.Configuration.OnosApiServer.PollInterval == 0 {
		ControllerConfig.Configuration.OnosApiServer.PollInterval = 30
	}

	logger.ControllerLog.Infoln("Ons Api Server Endpoint:")
	ControllerConfig.Configuration.OnosApiServer.Addr = strings.TrimSpace(ControllerConfig.Configuration.OnosApiServer.Addr)
	logger.ControllerLog.Infoln("Address ", ControllerConfig.Configuration.OnosApiServer.Addr)
	logger.ControllerLog.Infoln("Port ", ControllerConfig.Configuration.OnosApiServer.Port)
	logger.ControllerLog.Infoln("PollInterval ", ControllerConfig.Configuration.OnosApiServer.PollInterval)

	logger.ControllerLog.Infoln("Roc Endpoint:")
	ControllerConfig.Configuration.RocEndPoint.Addr = strings.TrimSpace(ControllerConfig.Configuration.RocEndPoint.Addr)
	logger.ControllerLog.Infoln("Address ", ControllerConfig.Configuration.RocEndPoint.Addr)
	logger.ControllerLog.Infoln("Port ", ControllerConfig.Configuration.RocEndPoint.Port)

	/*logger.ControllerLog.Infoln("Metric Func Endpoint:")
	ControllerConfig.Configuration.MetricFuncEndPoint.Addr = strings.TrimSpace(ControllerConfig.Configuration.MetricFuncEndPoint.Addr)
	logger.ControllerLog.Infoln("Address ", ControllerConfig.Configuration.MetricFuncEndPoint.Addr)
	logger.ControllerLog.Infoln("Port ", ControllerConfig.Configuration.MetricFuncEndPoint.Port)*/
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

func sendHttpReqMsg(req *http.Request) (*http.Response, error) {
	//Keep sending request to Http server until response is success
	var retries uint = 0
	var body []byte
	if req.Body != nil {
		body, _ = ioutil.ReadAll(req.Body)
	}
	for {
		cloneReq := req.Clone(context.Background())
		req.Body = ioutil.NopCloser(bytes.NewReader(body))
		cloneReq.Body = ioutil.NopCloser(bytes.NewReader(body))
		rsp, err := client.Do(cloneReq)
		retries += 1
		if err != nil {
			nextInterval := getNextBackoffInterval(retries, 2)
			logger.ControllerLog.Warningf("http req send error [%v], retrying after %v sec...", err.Error(), nextInterval)
			time.Sleep(time.Second * time.Duration(nextInterval))
			continue
		}

		if rsp.StatusCode == http.StatusAccepted ||
			rsp.StatusCode == http.StatusOK || rsp.StatusCode == http.StatusNoContent ||
			rsp.StatusCode == http.StatusCreated {
			logger.ControllerLog.Infoln("Get config from peer success")
			req.Body.Close()
			return rsp, nil
		} else {
			nextInterval := getNextBackoffInterval(retries, 2)
			logger.ControllerLog.Warningf("http rsp error [%v], retrying after [%v] sec...", http.StatusText(rsp.StatusCode), nextInterval)
			rsp.Body.Close()
			time.Sleep(time.Second * time.Duration(nextInterval))
		}
	}
}

func (onosClient *OnosService) GetRogueIPs(rogueIPChannel chan RogueIPs) {

	onosServerApi := onosClient.OnosServiceUrl + "/api/v1"
	req, err := http.NewRequest(http.MethodGet, onosServerApi, nil)
	if err != nil {
		logger.ControllerLog.Errorln("An Error Occured ", err)
		return
	}

	for {
		req.Header.Set("Content-Type", "application/json; charset=utf-8")

		rsp, httpErr := sendHttpReqMsg(req)
		if httpErr != nil {
			logger.ControllerLog.Errorf("Get Message [%v] returned error [%v] ", onosServerApi, err.Error())
			time.Sleep(10 * time.Second)
			continue
		}

		var rogueIPs RogueIPs
		if rsp != nil && rsp.Body != nil {
			json.NewDecoder(rsp.Body).Decode(&rogueIPs)
			logger.ControllerLog.Infoln("received rogueIPs from Onos App: ", rogueIPs)
			//writing rogueIPs into channel
			rogueIPChannel <- rogueIPs
		}
		time.Sleep(time.Duration(onosClient.PollInterval) * time.Second)
	}
}

/*func (metricClient *MetricFuncService) GetTargets(ipaddress string) (names []Targets) {
	metricApi := metricClient.MetricServiceUrl + "/nmetric-func/v1/subscriber/"
	req, err := http.NewRequest(http.MethodGet, rocTargetsApi, nil)
	if err != nil {
		fmt.Printf("An Error Occured %v", err)
		return
	}
	rsp, httpErr := sendHttpReqMsg(req)
	if httpErr != nil {
		log.Printf("Get Message [%v] returned error [%v] ", rocTargetsApi, err.Error())
	}

	if rsp != nil && rsp.Body != nil {
		json.NewDecoder(rsp.Body).Decode(&names)
		log.Printf("Targets received from RoC: %v", names)
	}
	return
}*/

func (rocClient *RocService) GetTargets() (names []Targets) {
	rocTargetsApi := rocClient.RocServiceUrl + "/aether-roc-api/targets"
	req, err := http.NewRequest(http.MethodGet, rocTargetsApi, nil)
	if err != nil {
		logger.ControllerLog.Errorf("An Error Occured %v", err)
		return
	}
	rsp, httpErr := sendHttpReqMsg(req)
	if httpErr != nil {
		logger.ControllerLog.Errorf("Get Message [%v] returned error [%v] ", rocTargetsApi, err.Error())
	}

	if rsp != nil && rsp.Body != nil {
		json.NewDecoder(rsp.Body).Decode(&names)
		logger.ControllerLog.Infoln("Targets received from RoC: ", names)
	}
	return
}

func (rocClient *RocService) DisableSimcard(targets []Targets, imsi string) {
	for _, target := range targets {
		rocSiteApi := rocClient.RocServiceUrl + "/aether-roc-api/aether/v2.1.x/" + target.EnterpriseId + "/site"
		req, err := http.NewRequest(http.MethodGet, rocSiteApi, nil)
		if err != nil {
			fmt.Printf("An Error Occured %v", err)
			return
		}
		rsp, httpErr := sendHttpReqMsg(req)
		if httpErr != nil {
			logger.ControllerLog.Errorf("Get Message [%v] returned error [%v] ", rocSiteApi, err.Error())
		}
		var siteInfo []SiteInfo
		if rsp != nil && rsp.Body != nil {
			json.NewDecoder(rsp.Body).Decode(&siteInfo)
			b, _ := io.ReadAll(rsp.Body)
			logger.ControllerLog.Infof("SimDetails Received from RoC: %s\n", string(b))
		}

		var rocDisableSimCard *SimCard
		for _, siteInfo := range siteInfo {
			for _, simCard := range siteInfo.SimCardDetails {
				if strings.HasPrefix(imsi, "imsi-") {
					imsi = imsi[5:]
				}
				if simCard.Imsi == imsi {
					logger.ControllerLog.Infof("SimCard %v Details Found in site [%v]\n", imsi, siteInfo.SiteId)
					rocDisableSimCard = &simCard
					break
				}
			}
			if rocDisableSimCard != nil {
				rocDisableImsiApi := rocSiteApi + "/" + siteInfo.SiteId + "/sim-card/" + rocDisableSimCard.SimId
				var val bool
				rocDisableSimCard.Enable = &val
				b, err := json.Marshal(&rocDisableSimCard)
				reqMsgBody := bytes.NewBuffer(b)
				fmt.Println("Rest API to disable IMSI: ", rocDisableImsiApi)
				fmt.Println("Post Msg Body:", reqMsgBody)

				req, err := http.NewRequest(http.MethodPost, rocDisableImsiApi, reqMsgBody)
				req.Header.Set("Content-Type", "application/json; charset=utf-8")
				_, httpErr := sendHttpReqMsg(req)
				if httpErr != nil {
					logger.ControllerLog.Errorf("Post Message [%v] returned error [%v] ", rocDisableImsiApi, err.Error())
				}
				break
			}
		}
	}
}

func RogueIPHandler(rogueIPChannel chan RogueIPs) {
	rocClient := RocService{
		RocServiceUrl: "http://" + ControllerConfig.Configuration.RocEndPoint.Addr + ":" + strconv.Itoa(ControllerConfig.Configuration.RocEndPoint.Port),
	}
	/*metricFuncClient := MetricService{
		MetricServiceUrl: "http://" + ControllerConfig.Configuration.MetricFuncEndPoint.Addr + ":" + strconv.Itoa(ControllerConfig.Configuration.MetricFuncEndPoint.Port),
	}*/

	for rogueIPs := range rogueIPChannel {

		for _, ipaddr := range rogueIPs.IpAddresses {
			// get IP to imsi mapping from metricfunc
			subscriberInfo, err := metricdata.GetSubscriberImsiFromIpAddr(ipaddr)
			if err != nil {
				logger.ControllerLog.Infoln("Subscriber Details doesn not exist with imsi ", err)
				continue
			}
			logger.ControllerLog.Infoln("Subscriber Imsi [%v] of the IP: [%v]", subscriberInfo.Imsi, ipaddr)
			//get enterprises or targets from ROC
			targets := rocClient.GetTargets()

			// get siteinfo from ROC
			//rocClient.DisableSimcard(targets, "208930100007490")
			rocClient.DisableSimcard(targets, subscriberInfo.Imsi)
		}
	}
}

/*func main() {
	rogueIpChan := make(chan RogueIPs, 100)
	InitConfigFactory()
	onosClient := OnosService{
		OnosServiceUrl: "http://" + ControllerConfig.Configuration.OnosApiServer.Addr + ":" +
			strconv.Itoa(ControllerConfig.Configuration.OnosApiServer.Port),
	}
	go onosClient.GetRogueIPs(rogueIpChan)
	go RogueIPHandler(rogueIpChan)

	select {}
}*/
