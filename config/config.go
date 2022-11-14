// SPDX-FileCopyrightText: 2022-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package config

type Config struct {
	Info          *Info          `yaml:"info"`
	Logger        *Logger        `yaml:"logger"`
	Configuration *Configuration `yaml:"configuration"`
}

type Info struct {
	Version     string `yaml:"version,omitempty"`
	Description string `yaml:"description,omitempty"`
}

type Logger struct {
	LogLevel string `yaml:"level,omitempty"`
}

type Configuration struct {
	NfStreams        []NFStream       `yaml:"nfStreams,omitempty"`
	AnalyticsStream  *AnalyticsStream `yaml:"analyticsStream,omitempty"`
	ApiServer        ServerAddr       `yaml:"apiServer,omitempty"`
	PrometheusServer ServerAddr       `yaml:"prometheusServer,omitempty"`
}

type ServerAddr struct {
	Addr string `yaml:"addr,omitempty"` // IP used to run the server in the node.
	Port int    `yaml:"port,omitempty"`
}

type NFStream struct {
	Urls  []string `yaml:"urls,omitempty"`
	Topic Topic    `yaml:"topic,omitempty"`
}

type Topic struct {
	TopicName   string `yaml:"topicName,omitempty"`
	TopicGroups string `yaml:"topicGroup,omitempty"`
}

type Groups struct {
	Analytics  string `yaml:"analytics,omitempty"`
	MongoDB    string `yaml:"mongodb,omitempty"`
	RestApis   string `yaml:"restapi,omitempty"`
	Prometheus string `yaml:"prometheus,omitempty"`
}

type AnalyticsStream struct {
	Enable    bool     `yaml:"enable,omitempty"`
	Urls      []string `yaml:"urls,omitempty"`
	TopicName string   `yaml:"topicName,omitempty"`
}
