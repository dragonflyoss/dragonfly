# Copyright The Dragonfly Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

PROJECT_NAME := "d7y.io/dragonfly/v2"
PKG := "$(PROJECT_NAME)"
PKG_LIST := $(shell go list ${PKG}/... | grep -v /vendor/ | grep -v '\(/test/\)')
GIT_COMMIT := $(shell git rev-parse --verify HEAD --short=7)
GIT_COMMIT_LONG := $(shell git rev-parse --verify HEAD)

all: help

# Prepare required folders for build.
build-dirs:
	@mkdir -p ./bin
.PHONY: build-dirs

# Build dragonfly.
docker-build: docker-build-scheduler docker-build-manager
	@echo "Build image done."
.PHONY: docker-build

# Push dragonfly images.
docker-push: docker-push-scheduler docker-push-manager
	@echo "Push image done."
.PHONY: docker-push

# Build scheduler image.
docker-build-scheduler:
	@echo "Begin to use docker build scheduler image."
	./hack/docker-build.sh scheduler
.PHONY: docker-build-scheduler

# Build manager image.
docker-build-manager:
	@echo "Begin to use docker build manager image."
	./hack/docker-build.sh manager
.PHONY: docker-build-manager

# Push scheduler image.
docker-push-scheduler: docker-build-scheduler
	@echo "Begin to push scheduler docker image."
	./hack/docker-push.sh scheduler
.PHONY: docker-push-scheduler

# Push manager image.
docker-push-manager: docker-build-manager
	@echo "Begin to push manager docker image."
	./hack/docker-push.sh manager
.PHONY: docker-push-manager

# Build dragonfly.
build: build-manager build-scheduler
.PHONY: build

# Build scheduler.
build-scheduler: build-dirs
	@echo "Begin to build scheduler."
	./hack/build.sh scheduler
.PHONY: build-scheduler

# Build manager.
build-manager: build-dirs build-manager-console
	@echo "Begin to build manager."
	make build-manager-server
.PHONY: build-manager

# Build manager server.
build-manager-server: build-dirs
	@echo "Begin to build manager server."
	./hack/build.sh manager
.PHONY: build-manager

# Build manager console.
build-manager-console: build-dirs
	@echo "Begin to build manager console."
	./hack/build.sh manager-console
.PHONY: build-manager-console

# Install scheduler.
install-scheduler:
	@echo "Begin to install scheduler."
	./hack/install.sh install scheduler
.PHONY: install-scheduler

# Install manager.
install-manager:
	@echo "Begin to install manager."
	./hack/install.sh install manager
.PHONY: install-manager

# Run unittests.
test:
	@go test -v -race -short ${PKG_LIST}
.PHONY: test

# Run tests with coverage.
test-coverage:
	@go test -v -race -short ${PKG_LIST} -coverprofile cover.out -covermode=atomic
	@cat cover.out >> coverage.txt
.PHONY: test-coverage

# Run github actions E2E tests with coverage.
actions-e2e-test-coverage:
	@ginkgo -v -r --race --fail-fast --cover --trace --show-node-events test/e2e
	@cat coverprofile.out >> coverage.txt
.PHONY: actions-e2e-test-coverage

# Run E2E tests.
e2e-test:
	@ginkgo -v -r --race --fail-fast --cover --trace --show-node-events test/e2e
.PHONY: e2e-test

# Run E2E tests with coverage.
e2e-test-coverage:
	@ginkgo -v -r --race --fail-fast --cover --trace --show-node-events test/e2e
	@cat coverprofile.out >> coverage.txt
.PHONY: e2e-test-coverage

# Clean E2E tests.
clean-e2e-test: 
	@echo "cleaning log file."
	@rm -rf test/e2e/*.log
.PHONY: clean-e2e-test

# Run code lint.
lint: markdownlint
	@echo "Begin to golangci-lint."
	@golangci-lint run
.PHONY: lint

# Run markdown lint.
markdownlint:
	@echo "Begin to markdownlint."
	@./hack/markdownlint.sh
.PHONY: markdownlint

# Run go generate.
generate:
	@go generate ${PKG_LIST}
.PHONY: generate

# Generate swagger files.
swag:
	@swag init --parseDependency --parseInternal -g cmd/manager/main.go -o api/manager

# Generate changelog.
changelog:
	@git-chglog -o CHANGELOG.md
.PHONY: changelog

# Clear compiled files.
clean:
	@go clean
	@rm -rf bin .go .cache
.PHONY: clean

fmt:
	@echo "Begin to go fmt."
	@go fmt ${PKG_LIST}
.PHONY: fmt

vet:
	@echo "Begin to go vet."
	@go vet ${PKG_LIST}
.PHONY: vet

precheck: fmt vet lint test
	@echo "All checks passed."
.PHONY: precheck

help: 
	@echo "make build-dirs                     prepare required folders for build"
	@echo "make docker-build                   build dragonfly image"
	@echo "make docker-push                    push dragonfly image"
	@echo "make docker-build-scheduler         build scheduler image"
	@echo "make docker-build-manager           build manager image"
	@echo "make docker-push-scheduler          push scheduler image"
	@echo "make docker-push-manager            push manager image"
	@echo "make build                          build dragonfly"
	@echo "make build-scheduler                build scheduler"
	@echo "make build-manager                  build manager"
	@echo "make build-manager-server           build manager server"
	@echo "make build-manager-console          build manager console"
	@echo "make install-scheduler              install scheduler"
	@echo "make install-manager                install manager"
	@echo "make test                           run unit tests"
	@echo "make test-coverage                  run tests with coverage"
	@echo "make actions-e2e-test-coverage      run github actons E2E tests with coverage"
	@echo "make e2e-test                       run e2e tests"
	@echo "make e2e-test-coverage              run e2e tests with coverage"
	@echo "make clean-e2e-test                 clean e2e tests"
	@echo "make lint                           run code lint"
	@echo "make markdownlint                   run markdown lint"
	@echo "make generate                       run go generate"
	@echo "make swag                           generate swagger api docs"
	@echo "make changelog                      generate CHANGELOG.md"
	@echo "make clean                          clean"
	@echo "make fmt                            run go fmt"
	@echo "make vet                            run go vet"
	@echo "make precheck                       run fmt, vet, lint, and test"
