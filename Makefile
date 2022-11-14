# SPDX-FileCopyrightText: 2022-present Intel Corporation
#
# SPDX-License-Identifier: Apache-2.0
#

PROJECT_NAME             := metricfunc
VERSION                  ?= $(shell cat ./VERSION)

## Docker related
DOCKER_REGISTRY          ?=
DOCKER_REPOSITORY        ?=
DOCKER_TAG               ?= ${VERSION}
DOCKER_IMAGENAME         := ${DOCKER_REGISTRY}${DOCKER_REPOSITORY}${PROJECT_NAME}:${DOCKER_TAG}
DOCKER_BUILDKIT          ?= 1
DOCKER_BUILD_ARGS        ?=

DOCKER_LABEL_BUILD_DATE  ?= $(shell date -u "+%Y-%m-%dT%H:%M:%SZ")

DOCKER_TARGETS           ?= metricfunc

# https://docs.docker.com/engine/reference/commandline/build/#specifying-target-build-stage---target

docker-build:
	@go mod vendor
	for target in $(DOCKER_TARGETS); do \
		DOCKER_BUILDKIT=$(DOCKER_BUILDKIT) buildx build --platform linux/amd64  $(DOCKER_BUILD_ARGS) \
			--target $$target \
			--tag ${DOCKER_REGISTRY}${DOCKER_REPOSITORY}5gc-$$target:${DOCKER_TAG} \
			--build-arg org_label_schema_version="${DOCKER_VERSION}" \
			--build-arg org_label_schema_vcs_url="${DOCKER_LABEL_VCS_URL}" \
			--build-arg org_label_schema_vcs_ref="${DOCKER_LABEL_VCS_REF}" \
			--build-arg org_label_schema_build_date="${DOCKER_LABEL_BUILD_DATE}" \
			--build-arg org_opencord_vcs_commit_date="${DOCKER_LABEL_COMMIT_DATE}" \
			. \
			|| exit 1; \
	done
	rm -rf vendor

docker-push:
	for target in $(DOCKER_TARGETS); do \
		docker push ${DOCKER_REGISTRY}${DOCKER_REPOSITORY}$$target:${DOCKER_TAG}; \
	done


.PHONY: docker-build docker-push


