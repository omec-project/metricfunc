// SPDX-FileCopyrightText: 2022-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package apiserver

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/omec-project/metricfunc/controller"
	"github.com/omec-project/metricfunc/internal/metricdata"
	"github.com/omec-project/metricfunc/logger"
	"github.com/omec-project/openapi"
)

func GetSubscriberSummary(c *gin.Context) {

	subId := c.Params.ByName("imsi")
	sub, _ := metricdata.GetSubscriber(subId)
	if sub != nil {
		resBody, err := openapi.Serialize(sub, "application/json")

		if err != nil {
			logger.ApiSrvLog.Errorf("json Marshal error %s", err.Error())
		}

		c.Writer.Write(resBody)
		c.Status(http.StatusOK)
		return
	}

	logger.ApiSrvLog.Errorf("subscriber data not found, imsi [%s] ", subId)
	c.JSON(http.StatusNotFound, gin.H{})
}

func GetSubscriberAll(c *gin.Context) {

	subs := metricdata.GetSubscriberAll()
	if len(subs) != 0 {
		resBody, err := openapi.Serialize(subs, "application/json")

		if err != nil {
			logger.ApiSrvLog.Errorf("json Marshal error %s", err.Error())
		}

		c.Writer.Write(resBody)
		c.Status(http.StatusOK)
		return
	}

	logger.ApiSrvLog.Errorf("no subscriber data not found ")
	c.JSON(http.StatusNotFound, gin.H{})
}

func GetNfStatus(c *gin.Context) {
	nfType := c.Params.ByName("type")

	nfs := metricdata.GetNfStatusbyNfType(nfType)
	if len(nfs) != 0 {
		resBody, err := openapi.Serialize(nfs, "application/json")

		if err != nil {
			logger.ApiSrvLog.Errorf("json Marshal error %s", err.Error())
		}

		c.Writer.Write(resBody)
		c.Status(http.StatusOK)
		return
	}
	logger.ApiSrvLog.Errorf("no nfs data not found ")
	c.JSON(http.StatusNotFound, gin.H{})
}

func GetNfStatusAll(c *gin.Context) {
	nfs := metricdata.GetNfStatusAll()

	if len(nfs) != 0 {
		resBody, err := openapi.Serialize(nfs, "application/json")

		if err != nil {
			logger.ApiSrvLog.Errorf("json Marshal error %s", err.Error())
		}
		c.Writer.Write(resBody)
		c.Status(http.StatusOK)
		return
	}
	logger.ApiSrvLog.Errorf("no nfs data not found ")
	c.JSON(http.StatusNotFound, gin.H{})
}

// Gives summary stats for any service
func GetNfServiceStatsSummary(c *gin.Context) {

}

// Gives detail stats of any service
func GetNfServiceStatsDetail(c *gin.Context) {
	nfType := c.Params.ByName("type")

	if svcStats, err := metricdata.GetNfServiceStatsDetail(nfType); err == nil {
		resBody, err := openapi.Serialize(svcStats, "application/json")
		if err != nil {
			logger.ApiSrvLog.Errorf("json Marshal error %s", err.Error())
		}
		c.Writer.Write(resBody)
		c.Status(http.StatusOK)
		return
	}
	logger.ApiSrvLog.Errorf("no nf service statistics data not found ")
	c.JSON(http.StatusNotFound, gin.H{})
}

// Gives summary of all services
func GetNfServiceStatsAll(c *gin.Context) {
}

func PushTestIPs(c *gin.Context) {
	requestBody, err := c.GetRawData()
	if err != nil {
		logger.ApiSrvLog.Errorf("get requestbody error %s", err.Error())
		return
	}
	var rogueIPs controller.RogueIPs
	json.Unmarshal(requestBody, &rogueIPs)

	logger.ApiSrvLog.Infoln("Test RogueIPs: ", rogueIPs)
	controller.RogueChannel <- rogueIPs
}
