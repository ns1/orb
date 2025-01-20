# Copyright (c) Mainflux
# SPDX-License-Identifier: Apache-2.0

# Adapted for Orb project, modifications licensed under MPL v. 2.0:
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/. */
include docker/.env

# expects to be set as env var
PRODUCTION_AGENT_REF_TAG ?= latest
PRODUCTION_AGENT_DEBUG_REF_TAG ?= latest-debug
REF_TAG ?= develop
DEBUG_REF_TAG ?= develop-debug
PKTVISOR_TAG ?= latest-develop
PKTVISOR_DEBUG_TAG ?= latest-develop-debug
DOCKER_IMAGE_NAME_PREFIX ?= orb
DOCKERHUB_REPO = orbcommunity
ORB_DOCKERHUB_REPO = orbcommunity
BUILD_DIR = build
SERVICES = fleet policies sinks sinker migrate maestro
DOCKERS = $(addprefix docker_,$(SERVICES))
DOCKERS_DEV = $(addprefix docker_dev_,$(SERVICES))
CGO_ENABLED ?= 0
GOARCH ?= $(shell dpkg-architecture -q DEB_BUILD_ARCH)
GOOS ?= $(shell dpkg-architecture -q DEB_TARGET_ARCH_OS)
ORB_VERSION = $(shell cat VERSION)
COMMIT_HASH = $(shell git rev-parse --short HEAD)
OTEL_COLLECTOR_CONTRIB_VERSION ?= 0.91.0
OTEL_CONTRIB_URL ?= "https://github.com/open-telemetry/opentelemetry-collector-releases/releases/download/v$(OTEL_COLLECTOR_CONTRIB_VERSION)/otelcol-contrib_$(OTEL_COLLECTOR_CONTRIB_VERSION)_$(GOOS)_$(GOARCH).tar.gz"

define compile_service
    echo "ORB_VERSION: $(ORB_VERSION)"
	CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GOOS) GOARCH=$(GOARCH) GOARM=$(GOARM) go build -mod=mod -ldflags "-extldflags "-static" -X 'github.com/orb-community/orb/buildinfo.version=$(ORB_VERSION)'" -o ${BUILD_DIR}/$(DOCKER_IMAGE_NAME_PREFIX)-$(1) cmd/$(1)/main.go
endef

define compile_service_linux
	$(eval svc=$(subst docker_dev_,,$(1)))
    echo "ORB_VERSION: $(ORB_VERSION)-$(COMMIT_HASH)"
	CGO_ENABLED=$(CGO_ENABLED) GOOS=linux GOARCH=$(GOARCH) GOARM=$(GOARM) go build -mod=mod -ldflags "-extldflags "-static" -X 'github.com/orb-community/orb/buildinfo.version=$(ORB_VERSION)-$(COMMIT_HASH)'" -o ${BUILD_DIR}/$(DOCKER_IMAGE_NAME_PREFIX)-$(svc) cmd/$(svc)/main.go
endef

define run_test
	 go test -mod=mod -short -race -count 1 -tags test $(shell go list ./... | grep -v 'cmd' | grep '$(SERVICE)')
endef

define run_test_coverage
	 go test -mod=mod -short -race -count 1 -tags test -cover -coverprofile=coverage.out -covermode=atomic $(shell go list ./... | grep -v 'cmd' | grep '$(SERVICE)')
endef

define make_docker
	$(eval SERVICE=$(shell [ -z "$(SERVICE)" ] && echo $(subst docker_,,$(1)) || echo $(SERVICE)))
	docker build \
		--no-cache \
		--build-arg SVC=$(SERVICE) \
		--build-arg GOARCH=$(GOARCH) \
		--build-arg GOARM=$(GOARM) \
		--tag=$(ORB_DOCKERHUB_REPO)/$(DOCKER_IMAGE_NAME_PREFIX)-$(SERVICE):$(REF_TAG) \
		--tag=$(ORB_DOCKERHUB_REPO)/$(DOCKER_IMAGE_NAME_PREFIX)-$(SERVICE):$(ORB_VERSION) \
		--tag=$(ORB_DOCKERHUB_REPO)/$(DOCKER_IMAGE_NAME_PREFIX)-$(SERVICE):$(ORB_VERSION)-$(COMMIT_HASH) \
		-f docker/Dockerfile .
	$(eval SERVICE="")
endef
define make_docker_debug
	$(eval SERVICE=$(shell [ -z "$(SERVICE)" ] && echo $(subst docker_,,$(1)) || echo $(SERVICE)))
	docker build \
		--no-cache \
		--build-arg SVC=$(SERVICE) \
		--build-arg GOARCH=$(GOARCH) \
		--build-arg GOARM=$(GOARM) \
		--tag=$(ORB_DOCKERHUB_REPO)/$(DOCKER_IMAGE_NAME_PREFIX)-$(SERVICE):$(REF_TAG) \
		--tag=$(ORB_DOCKERHUB_REPO)/$(DOCKER_IMAGE_NAME_PREFIX)-$(SERVICE):$(ORB_VERSION) \
		--tag=$(ORB_DOCKERHUB_REPO)/$(DOCKER_IMAGE_NAME_PREFIX)-$(SERVICE):$(ORB_VERSION)-$(COMMIT_HASH) \
		-f docker/Dockerfile.debug .
	$(eval SERVICE="")
endef

define make_docker_dev
	$(eval svc=$(shell [ -z "$(SERVICE)" ] && echo $(subst docker_dev_,,$(1)) || echo $(svc)))
	docker build \
		--no-cache \
		--build-arg SVC=$(svc) \
		--tag=$(ORB_DOCKERHUB_REPO)/$(DOCKER_IMAGE_NAME_PREFIX)-$(svc):$(REF_TAG) \
		--tag=$(ORB_DOCKERHUB_REPO)/$(DOCKER_IMAGE_NAME_PREFIX)-$(svc):$(ORB_VERSION) \
		--tag=$(ORB_DOCKERHUB_REPO)/$(DOCKER_IMAGE_NAME_PREFIX)-$(svc):$(ORB_VERSION)-$(COMMIT_HASH) \
		-f docker/Dockerfile.dev ./build
	$(eval svc="")
endef

all: platform

.PHONY: all $(SERVICES) dockers dockers_dev ui services agent agent_bin

clean:
	rm -rf ${BUILD_DIR}

cleandocker:
	# Stops containers and removes containers, networks, volumes, and images created by up
#	docker-compose -f docker/docker-compose.yml down --rmi all -v --remove-orphans
	docker-compose -f docker/docker-compose.yml down -v --remove-orphans

ifdef pv
	# Remove unused volumes
	docker volume ls -f name=$(DOCKER_IMAGE_NAME_PREFIX) -f dangling=true -q | xargs -r docker volume rm
endif

test:
	go test -mod=mod -short -race -count 1 -tags test $(shell go list ./... | grep -v 'cmd')

run_test_service: test_service $(2)

run_test_service_cov: test_service_cov $(2)

test_service:
	$(call run_test,$(@))

test_service_cov:
	$(call run_test_coverage,$(@))

proto:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative policies/pb/policies.proto
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative fleet/pb/fleet.proto
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative sinks/pb/sinks.proto

$(SERVICES):
	$(call compile_service,$(@))

$(DOCKERS):
	$(call make_docker,$(@),$(GOARCH))

$(DOCKERS_DEV):
	$(call compile_service_linux,$(@))
	$(call make_docker_dev,$(@))

services: $(SERVICES)
dockers: $(DOCKERS)
dockers_dev: $(DOCKERS_DEV)

build_docker:
	$(call make_docker,$(@),$(GOARCH))

build_docker_debug:
	$(call make_docker_debug,$(@),$(GOARCH))

# install tools for kind

install-kind:
	cd /tmp && \
	curl -Lo ./kind https://kind.sigs.k8s.io/dl/v0.17.0/kind-linux-amd64 && \
	chmod +x ./kind && \
	sudo mv ./kind /usr/local/bin/kind

install-helm:
	cd /tmp
	curl https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash

install-kubectl:
	cd /tmp && \
	curl -LO "https://dl.k8s.io/release/v1.27.4/bin/linux/amd64/kubectl" && \
	chmod a+x ./kubectl && \
	sudo mv ./kubectl /usr/local/bin/kubectl

install-docker:
	cd /tmp
	curl -fsSL https://get.docker.com -o get-docker.sh
	sh ./get-docker.sh

install-k9s:
	cd /tmp && \
	wget https://github.com/derailed/k9s/releases/download/v0.26.7/k9s_Linux_x86_64.tar.gz && \
	tar -xvzf k9s_Linux_x86_64.tar.gz && \
	sudo install -o root -g root -m 0755 k9s /usr/local/bin/k9s


# kind commands

prepare-helm:
	cd ./kind/ && \
	helm repo add jaegertracing https://jaegertracing.github.io/helm-charts && \
	helm repo add orb-community https://orb-community.github.io/orb-helm/ && \
	helm dependency build

kind-create-all: kind-create-cluster kind-install-orb

kind-upgrade-all: kind-load-images kind-upgrade-orb

kind-create-cluster:
	kind create cluster --image kindest/node:v1.24.0 --config=./kind/config.yaml

kind-delete-cluster:
	kind delete cluster

kind-load-images:
	kind load docker-image orbcommunity/orb-fleet:develop
	kind load docker-image orbcommunity/orb-policies:develop
	kind load docker-image orbcommunity/orb-sinks:develop
	kind load docker-image orbcommunity/orb-sinker:develop
	kind load docker-image orbcommunity/orb-migrate:develop
	kind load docker-image orbcommunity/orb-maestro:develop
	kind load docker-image orbcommunity/orb-ui:develop

kind-install-orb:
	kubectl create namespace orb
	kubectl create namespace otelcollectors
	kubectl create secret generic orb-auth-service --from-literal=jwtSecret=MY_SECRET -n orb
	kubectl create secret generic orb-user-service --from-literal=adminEmail=admin@kind.com --from-literal=adminPassword=pass123456 -n orb
	kubectl create secret generic orb-sinks-encryption-key --from-literal=key=MY_SINKS_SECRET -n orb
	helm install -n orb kind-orb ./kind
	kubectl apply -f ./kind/nginx.yaml

kind-upgrade-orb:
	helm upgrade -n orb kind-orb ./kind
	kubectl rollout restart deployment -n orb

kind-delete-orb:
	kubectl delete -f ./kind/nginx.yaml
	helm delete -n orb kind-orb
	kubectl delete secret orb-user-service -n orb
	kubectl delete secret orb-auth-service -n orb
	kubectl delete namespace orb
	kubectl delete namespace otelcollectors

#

run: prepare-helm kind-create-all

stop: kind-delete-orb kind-delete-cluster

agent_bin:
	$(call compile_service_linux,agent)

agent:
	docker build --no-cache \
	  --build-arg GOARCH=$(GOARCH) \
	  --build-arg PKTVISOR_TAG=$(PKTVISOR_TAG) \
	  --tag=$(ORB_DOCKERHUB_REPO)/$(DOCKER_IMAGE_NAME_PREFIX)-agent:$(REF_TAG) \
	  --tag=$(ORB_DOCKERHUB_REPO)/$(DOCKER_IMAGE_NAME_PREFIX)-agent:$(ORB_VERSION) \
	  --tag=$(ORB_DOCKERHUB_REPO)/$(DOCKER_IMAGE_NAME_PREFIX)-agent:$(ORB_VERSION)-$(COMMIT_HASH) \
	  -f agent/docker/Dockerfile .
	  
agent_full:
	docker build --no-cache \
	  --build-arg GOARCH=$(GOARCH) \
	  --build-arg PKTVISOR_TAG=$(PKTVISOR_TAG) \
	  --build-arg ORB_TAG=${REF_TAG} \
	  --build-arg OTEL_TAG=${OTEL_COLLECTOR_CONTRIB_VERSION} \
	  --tag=$(ORB_DOCKERHUB_REPO)/$(DOCKER_IMAGE_NAME_PREFIX)-agent-full:$(REF_TAG) \
	  --tag=$(ORB_DOCKERHUB_REPO)/$(DOCKER_IMAGE_NAME_PREFIX)-agent-full:$(ORB_VERSION) \
	  --tag=$(ORB_DOCKERHUB_REPO)/$(DOCKER_IMAGE_NAME_PREFIX)-agent-full:$(ORB_VERSION)-$(COMMIT_HASH) \
	  -f agent/docker/Dockerfile.full .

agent_custom:
	docker build --no-cache \
	  --build-arg GOARCH=$(GOARCH) \
	  --build-arg PKTVISOR_TAG=$(PKTVISOR_TAG) \
	  --build-arg DIODE_TAG=$(DIODE_TAG) \
	  --build-arg ORB_TAG=${REF_TAG} \
	  --build-arg OTEL_TAG=latest \
	  --tag=$(ORB_DOCKERHUB_REPO)/$(DOCKER_IMAGE_NAME_PREFIX)-agent-custom:$(REF_TAG) \
	  --tag=$(ORB_DOCKERHUB_REPO)/$(DOCKER_IMAGE_NAME_PREFIX)-agent-custom:$(ORB_VERSION) \
	  --tag=$(ORB_DOCKERHUB_REPO)/$(DOCKER_IMAGE_NAME_PREFIX)-agent-custom:$(ORB_VERSION)-$(COMMIT_HASH) \
	  -f agent/docker/Dockerfile.custom .

agent_debug:
	docker build \
	  --build-arg PKTVISOR_TAG=$(PKTVISOR_DEBUG_TAG) \
	  --tag=$(DOCKERHUB_REPO)/$(DOCKER_IMAGE_NAME_PREFIX)-agent:$(DEBUG_REF_TAG) \
	  --tag=$(ORB_DOCKERHUB_REPO)/$(DOCKER_IMAGE_NAME_PREFIX)-agent:$(DEBUG_REF_TAG) \
	  -f agent/docker/Dockerfile .

agent_production:
	docker build \
	  --build-arg PKTVISOR_TAG=$(PKTVISOR_TAG) \
	  --tag=$(ORB_DOCKERHUB_REPO)/$(DOCKER_IMAGE_NAME_PREFIX)-agent:$(PRODUCTION_AGENT_REF_TAG) \
	  --tag=$(ORB_DOCKERHUB_REPO)/$(DOCKER_IMAGE_NAME_PREFIX)-agent:$(ORB_VERSION) \
	  --tag=$(ORB_DOCKERHUB_REPO)/$(DOCKER_IMAGE_NAME_PREFIX)-agent:$(ORB_VERSION)-$(COMMIT_HASH) \
	  -f agent/docker/Dockerfile .

agent_debug_production:
	docker build \
	  --build-arg PKTVISOR_TAG=$(PKTVISOR_DEBUG_TAG) \
	  --tag=$(ORB_DOCKERHUB_REPO)/$(DOCKER_IMAGE_NAME_PREFIX)-agent:$(PRODUCTION_AGENT_DEBUG_REF_TAG) \
	  -f agent/docker/Dockerfile .

test_ui:
	cd ui/ && yarn test

ui-modules:
	cd ui/ && docker build \
		--tag=$(ORB_DOCKERHUB_REPO)/orb-ui-modules:latest \
		--tag=$(ORB_DOCKERHUB_REPO)/orb-ui-modules:$(REF_TAG) \
		--tag=$(ORB_DOCKERHUB_REPO)/orb-ui-modules:$(ORB_VERSION)-$(COMMIT_HASH) \
		-f docker/Dockerfile.buildyarn .

ui:
	cd ui/ && docker build \
		--build-arg ENV=${ENVIRONMENT} \
		--tag=$(ORB_DOCKERHUB_REPO)/$(DOCKER_IMAGE_NAME_PREFIX)-ui:$(REF_TAG) \
		--tag=$(ORB_DOCKERHUB_REPO)/$(DOCKER_IMAGE_NAME_PREFIX)-ui:$(ORB_VERSION) \
		--tag=$(ORB_DOCKERHUB_REPO)/$(DOCKER_IMAGE_NAME_PREFIX)-ui:$(ORB_VERSION)-$(COMMIT_HASH) \
		-f docker/Dockerfile .

platform: dockers_dev agent ui

pull-latest-otel-collector-contrib:
	wget -O ./agent/backend/otel/otelcol_contrib.tar.gz $(OTEL_CONTRIB_URL)
	tar -xvf ./agent/backend/otel/otelcol_contrib.tar.gz -C ./agent/backend/otel/
	cp ./agent/backend/otel/otelcol-contrib .
	rm ./agent/backend/otel/otelcol_contrib.tar.gz
	rm ./agent/backend/otel/LICENSE
	rm ./agent/backend/otel/README.md
