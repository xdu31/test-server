REPO       := github.com/xdu31/test-server
BUILD_PATH := bin

CUR_DIR := $(shell dirname $(abspath $(lastword $(MAKEFILE_LIST))))

DOCKERFILE_PATH := $(CUR_DIR)/docker
BIN_DIR          = $(CUR_DIR)/bin

BUILDTOOL_IMAGE := golang:1.14.1-alpine3.11
GENTOOL_IMAGE := infoblox/atlas-gentool:latest

# configuration for image names
USERNAME       := "xdu31"
GIT_COMMIT     := $(shell git describe --always || echo pre-commit)
#GIT_COMMIT     := $(shell git describe --dirty=-unsupported --always || echo pre-commit)
IMAGE_VERSION  ?= $(GIT_COMMIT)
IMAGE_REGISTRY ?= xdu31

APP_SERVER     := server
PACKAGE_SERVER := $(REPO)/cmd/$(APP_SERVER)
PACKAGE_INTEGRATION_TEST := $(REPO)/integration

IMAGE_NAME              ?= test-server
DB_MIGRATION_IMAGE_NAME ?= test-server-db-migrate
# configuration for server binary and image
SERVER_BINARY      := $(BUILD_PATH)/$(APP_SERVER)
SERVER_IMAGE       := $(IMAGE_REGISTRY)/$(IMAGE_NAME)
DB_MIGRATION_IMAGE := $(IMAGE_REGISTRY)/$(DB_MIGRATION_IMAGE_NAME)
SERVER_DOCKERFILE  := $(DOCKERFILE_PATH)/Dockerfile
DB_MIGRATION_DOCKERFILE := $(DOCKERFILE_PATH)/Dockerfile.migrate

# configuration for building on host machine
GO_CACHE := $(BUILD_PATH)/go-cache
SRCROOT  := /go/src/$(REPO)

#TODO: Some error will occur 
#GO_BUILD_FLAGS ?= -pkgdir $(SRCROOT)/$(GO_CACHE) -i -v -ldflags '-w -s -X main.Version=$(GIT_COMMIT) -X main.GitCommit=$(GIT_COMMIT)'
GO_BUILD_FLAGS ?= -pkgdir $(SRCROOT)/$(GO_CACHE) -i -v -ldflags '-w -s'
GO_TEST_FLAGS  ?= -v -cover

DOCKER_RUNNER        := docker run --rm
DOCKER_RUNNER        := docker run
DOCKER_RUNNER        += -v $(CUR_DIR):$(SRCROOT)
DOCKER_RUNNER        += -w $(SRCROOT)

BUILDER               = $(DOCKER_RUNNER) -e CGO_ENABLED=0 -e GO111MODULE=off $(BUILDTOOL_IMAGE)
TESTER                = $(DOCKER_RUNNER) -e CGO_ENABLED=0 -e GO111MODULE=off $(BUILDTOOL_IMAGE)

GOBUILD  = $(BUILDER) go build $(GO_BUILD_FLAGS)
GOFMT    = $(BUILDER) go fmt
GORUN    = $(BUILDER) go run

DOCKER_GENTOOL_DEBUG_ENV  = --env=GRPC_TRACE=all
DOCKER_GENTOOL_DEBUG_ENV += --env=GRPC_VERBOSITY=DEBUG
DOCKER_GENTOOL = docker run $(DOCKER_GENTOOL_DEBUG_ENV) --rm -v $(CUR_DIR):$(SRCROOT) $(GENTOOL_IMAGE)
GENTOOL_PARAM_NOPREPROCESS  = --gorm_out=.
GENTOOL_PARAM_NOPREPROCESS += --go_out=plugins=grpc:.
GENTOOL_PARAM_NOPREPROCESS += --grpc-gateway_out=logtostderr=true,allow_delete_body=true:.
GENTOOL_PARAM_NOPREPROCESS += --swagger_out="atlas_patch=true,allow_delete_body=true:."

GENTOOL_PARAM_NOPREPROCESS += -I $(SRCROOT)/vendor
GENTOOL_PARAM_NOPREPROCESS += --validate_out="lang=go:."

GENTOOL_PARAM = $(GENTOOL_PARAM_NOPREPROCESS) --preprocess_out=.
GENTOOL = $(DOCKER_GENTOOL) $(GENTOOL_PARAM)

GENTOOL_NOPREPROCESS = $(DOCKER_GENTOOL) $(GENTOOL_PARAM_NOPREPROCESS)

GENTOOL_SWAGGER_PARAM  = --swagger_out="atlas_patch=true,allow_delete_body=true:."
GENTOOL_SWAGGER_PARAM += -I $(SRCROOT)/vendor
SWAGGER_GENTOOL = $(DOCKER_GENTOOL) $(GENTOOL_SWAGGER_PARAM)

default: build

.PHONY: vendor
vendor:
	GO111MODULE=on GOPROXY=direct GOSUMDB=off
	go mod vendor
	go mod download

.PHONY: gen-buf
gen-buf:
	$(GENTOOL_NOPREPROCESS) \
	$(REPO)/pkg/pb/service.proto

.PHONY: gen-buf-swagger
gen-buf-swagger:
	$(SWAGGER_GENTOOL) \
	$(REPO)/pkg/pb/service.proto

.PHONY: gen
gen: gen-buf gen-buf-swagger

.PHONY: bin
bin:
	@mkdir -p "$(BIN_DIR)"

.PHONY: fmt
fmt:
	@echo "Running 'go fmt ...'"
	$(GOFMT) -x "$(REPO)/cmd/..." \
		"$(REPO)/pkg/..."

.PHONY: build-server
build-server: bin
	$(GOBUILD) -o $(SERVER_BINARY) $(PACKAGE_SERVER)

.PHONY: test
test: fmt
	@go test $(GO_TEST_FLAGS) -tags=unit $(PACKAGE_SERVER)

.PHONY: build
build: fmt build-server

.PHONY: test-with-integration
test-with-integration: fmt
	@go test $(GO_TEST_FLAGS) -tags=integration $(PACKAGE_INTEGRATION_TEST)

.PHONY: image
image: server-image db-migration-image

server-image:
	@docker build -f $(SERVER_DOCKERFILE) -t $(SERVER_IMAGE):$(IMAGE_VERSION) .
	@docker tag $(SERVER_IMAGE):$(IMAGE_VERSION) $(SERVER_IMAGE):latest
#	TODO update my docker Daemon client API version to 1.25+
#	@docker image prune -f --filter label=stage=server-intermediate

db-migration-image:
	@docker build -f $(DB_MIGRATION_DOCKERFILE) -t $(DB_MIGRATION_IMAGE):$(IMAGE_VERSION) .
	@docker tag $(DB_MIGRATION_IMAGE):$(IMAGE_VERSION) $(DB_MIGRATION_IMAGE):latest
#	TODO update my docker Daemon client API version to 1.25+
#	@docker image prune -f --filter label=stage=server-intermediate

.PHONY: push
push:
	@docker push $(SERVER_IMAGE)
	@docker push $(DB_MIGRATION_IMAGE)

.PHONY: clean
clean:
	$(BUILDER) sh -c "test ! -d $(SRCROOT)/bin || chmod -R 777 $(SRCROOT)/bin"
	@rm -rf "$(BIN_DIR)"
	@docker rmi -f $(shell docker images -q $(SERVER_IMAGE)) || true
	@docker rmi -f $(shell docker images -q $(DB_MIGRATION_IMAGE)) || true

