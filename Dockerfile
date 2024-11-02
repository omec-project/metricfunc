# SPDX-FileCopyrightText: 2022-present Intel Corporation
#
# SPDX-License-Identifier: Apache-2.0
#

FROM golang:1.23.2-bookworm AS builder

WORKDIR $GOPATH/src/metricfunc
COPY . .
RUN make all

FROM alpine:3.20 AS metricfunc

LABEL maintainer="Aether SD-Core <dev@lists.aetherproject.org>" \
    description="Aether open source 5G Core Network" \
    version="Stage 3"

ARG DEBUG_TOOLS

# Install debug tools ~ 50MB (if DEBUG_TOOLS is set to true)
RUN if [ "$DEBUG_TOOLS" = "true" ]; then \
        apk update && apk add --no-cache -U vim strace net-tools curl netcat-openbsd bind-tools tcpdump; \
        fi

# Copy executable
COPY --from=builder /go/src/metricfunc/bin/* /usr/local/bin/.
