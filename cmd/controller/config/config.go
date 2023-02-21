// SPDX-FileCopyrightText: 2022-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package Config

type Config struct {
	Info          *Info          `yaml:"info"`
	Logger        *Logger        `yaml:"logger"`
	Configuration *Configuration `yaml:"configuration"`
}

type Info struct {
	Version     string `yaml:"version,omitempty"`
	Description string `yaml:"description,omitempty"`
	HttpVersion int    `yaml:"http-version,omitempty"`
}

type Logger struct {
	LogLevel string `yaml:"level,omitempty"`
}

type Configuration struct {
	OnosApiServer      ServerAddr `yaml:"onosApiServer,omitempty"`
	RocEndPoint        ServerAddr `yaml:"rocEndPoint,omitempty"`
	MetricFuncEndPoint ServerAddr `yaml:"metricFuncEndPoint,omitempty"`
}

type ServerAddr struct {
	Addr         string `yaml:"addr,omitempty"` // IP used to run the server in the node.
	Port         int    `yaml:"port,omitempty"`
	PollInterval int    `yaml:"pollInterval,omitempty"`
}

type Urls struct {
	Uri  string `yaml:"uri,omitempty"`
	Port int    `yaml:"port,omitempty"`
}
