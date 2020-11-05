# Copyright 2015 The Prometheus Authors
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

GO     ?= GO15VENDOREXPERIMENT=1 go
GOPATH := $(firstword $(subst :, ,$(shell $(GO) env GOPATH)))
GOARCH := $(shell $(GO) env GOARCH)
GOHOSTARCH := $(shell $(GO) env GOHOSTARCH)

PROMTOOL    ?= $(GOPATH)/bin/promtool
PROMU       ?= $(GOPATH)/bin/promu
STATICCHECK ?= $(GOPATH)/bin/staticcheck
pkgs         = $(shell $(GO) list ./... | grep -v /vendor/)

PREFIX                  ?= $(shell pwd)
BIN_DIR                 ?= $(shell pwd)
DOCKER_IMAGE_NAME       ?= opensips_exporter
DOCKER_IMAGE_TAG        ?= $(subst /,-,$(shell git rev-parse --abbrev-ref HEAD))
MACH                    ?= $(shell uname -m)
DOCKERFILE              ?= Dockerfile

NAME ?= opensips
OPENSIPS_DOCKER_TAG ?= 2.4.0

all: format vet staticcheck build

style:
	@echo ">> checking code style"
	@! gofmt -d $(shell find . -path ./vendor -prune -o -name '*.go' -print) | grep '^'

format:
	@echo ">> formatting code"
	@$(GO) fmt $(pkgs)

vet:
	@echo ">> vetting code"
	@$(GO) vet $(pkgs)

staticcheck: $(STATICCHECK)
	@echo ">> running staticcheck"
	@$(STATICCHECK) -ignore "$(STATICCHECK_IGNORE)" $(pkgs)

build: $(PROMU)
	@echo ">> building binaries"
	@$(PROMU) build --prefix $(PREFIX)

tarball: $(PROMU)
	@echo ">> building release tarball"
	@$(PROMU) tarball --prefix $(PREFIX) $(BIN_DIR)

docker:
	@echo ">> building docker image from $(DOCKERFILE)"
	@docker build --file $(DOCKERFILE) -t "$(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG)" .

docker-build:
	docker build \
		--tag="opensips/opensips:$(OPENSIPS_DOCKER_TAG)" \
		.

start:
	docker run -p 8888:8888 -d --name $(NAME) opensips/opensips:$(OPENSIPS_DOCKER_TAG) -v /home/ruben/go/src/github.com/VoIPGRID/opensips_exporter/:/tmp/


$(GOPATH)/bin/promtool promtool:
	@GOOS= GOARCH= $(GO) get -u github.com/prometheus/prometheus/cmd/promtool

$(GOPATH)/bin/promu promu:
	@GOOS= GOARCH= $(GO) get -u github.com/prometheus/promu

$(GOPATH)/bin/staticcheck:
	@GOOS= GOARCH= $(GO) get -u honnef.co/go/tools/cmd/staticcheck


.PHONY: all style format build test vet tarball docker promtool promu staticcheck

# Declaring the binaries at their default locations as PHONY targets is a hack
# to ensure the latest version is downloaded on every make execution.
# If this is not desired, copy/symlink these binaries to a different path and
# set the respective environment variables.
.PHONY: $(GOPATH)/bin/promtool $(GOPATH)/bin/promu $(GOPATH)/bin/staticcheck
