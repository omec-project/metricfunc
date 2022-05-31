# SPDX-FileCopyrightText: 2022-present Intel Corporation
#
# SPDX-License-Identifier: Apache-2.0
#

FROM golang:1.18.0-stretch AS builder

LABEL maintainer="ONF <omec-dev@opennetworking.org>"

RUN apt-get update
RUN apt-get -y install vim 
RUN cd $GOPATH/src && mkdir -p metricfunc
COPY . $GOPATH/src/metricfunc
RUN cd $GOPATH/src/metricfunc && go install cmd/metricfunc/metricfunc.go 
RUN cd $GOPATH/src/metricfunc && go install cmd/client/client.go 

FROM builder AS metricfunc
WORKDIR /metricfunc
RUN mkdir -p /metricfunc/bin
COPY --from=builder $GOPATH/bin/* /metricfunc/bin/
WORKDIR /metricfunc
