# SPDX-FileCopyrightText: 2022-present Intel Corporation
#
# SPDX-License-Identifier: Apache-2.0
#

FROM golang:1.21.6-bookworm AS builder

LABEL maintainer="ONF <omec-dev@opennetworking.org>"

RUN apt-get update && apt-get -y install vim


RUN cd $GOPATH/src && mkdir -p metricfunc
COPY . $GOPATH/src/metricfunc
RUN cd $GOPATH/src/metricfunc/cmd/metricfunc && CGO_ENABLED=0 go build -mod=mod

FROM alpine:3.16 as metricfunc

ARG DEBUG_TOOLS
# Install debug tools ~ 100MB (if DEBUG_TOOLS is set to true)
RUN apk update && apk add -U vim strace net-tools curl netcat-openbsd bind-tools bash tcpdump 

# Set working dir
WORKDIR /metricfunc
RUN mkdir -p /metricfunc/bin

# Copy executable
COPY --from=builder /go/src/metricfunc/cmd/metricfunc/metricfunc /metricfunc/bin/

#Image default directory
WORKDIR /metricfunc
