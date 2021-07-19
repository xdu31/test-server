REPO       := github.com/xdu31/test-server

BUILD_PATH              := bin
CUR_DIR                 := $(shell dirname $(abspath $(lastword $(MAKEFILE_LIST))))
DOCKERFILE_PATH         := $(CUR_DIR)/docker
BIN_DIR                 := $(CUR_DIR)/bin

BUILDTOOL_IMAGE         := golang:1.16.0-alpine3.13
GENTOOL_IMAGE           := infoblox/atlas-gentool:v22.0

GIT_COMMIT              := $(shell git describe --tag --dirty=-unsupported --always || echo pre-commit)

IMAGE_REGISTRY ?= xdu31
IMAGE_PREFIX ?= test-server
IMAGE_VERSION  ?= $(GIT_COMMIT)

SERVER_APP              := server
SERVER_BINARY           := $(BUILD_PATH)/$(SERVER_APP)
SERVER_PATH             := $(REPO)/cmd/$(SERVER_APP)
SERVER_IMAGE_NAME       := $(IMAGE_PREFIX)
SERVER_IMAGE_FULL       := $(IMAGE_REGISTRY)/$(SERVER_IMAGE_NAME)
SERVER_DOCKERFILE       := $(CURDIR)/cmd/$(SERVER_APP)/Dockerfile

DB_MIGRATION_IMAGE_NAME := $(SERVER_IMAGE_NAME)-db-migrate
DB_MIGRATION_IMAGE_FULL := $(IMAGE_REGISTRY)/$(DB_MIGRATION_IMAGE_NAME)
DB_MIGRATION_DOCKERFILE := $(CURDIR)/db/Dockerfile

GO_CACHE := $(BUILD_PATH)/go-cache
SRCROOT  := /go/src/$(REPO)

GO_BUILD_FLAGS ?= -pkgdir $(GO_CACHE) -i -v -ldflags '-w -s'
GO_TEST_FLAGS  ?= -v -cover

DOCKER_RUNNER        := docker run --rm
DOCKER_RUNNER        += -v $(CUR_DIR):$(SRCROOT)
DOCKER_RUNNER        += -w $(SRCROOT)

BUILDER               = $(DOCKER_RUNNER) -e CGO_ENABLED=0 -e GO111MODULE=off $(BUILDTOOL_IMAGE)
TESTER                = $(DOCKER_RUNNER) -e CGO_ENABLED=0 -e GO111MODULE=off $(BUILDTOOL_IMAGE)

NAMESPACE               ?= test-server
CHART                   ?= $(NAMESPACE)
CHART_FILE              := $(CHART)-$(IMAGE_VERSION).tgz
CHART_PATH              := $(PWD)/helm

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

GENTOOL_PARAM = $(GENTOOL_PARAM_NOPREPROCESS) --preprocess_out=. $(JSONSCHEMA)
GENTOOL = $(DOCKER_GENTOOL) $(GENTOOL_PARAM)

GENTOOL_NOPREPROCESS = $(DOCKER_GENTOOL) $(GENTOOL_PARAM_NOPREPROCESS)

SWAGGER_FILE            := pkg/server/pb/service.swagger.json
SWAGGER_EXCLUDE_LIST    := test_data

HELM_IMAGE              ?= infoblox/helm:2.14.3-1
HELM_GENERATOR          ?= docker run --rm -v $(CHART_PATH):/repo -e AWS_REGION=${AWS_REGION} -e AWS_ACCESS_KEY_ID=${AWS_ACCESS_KEY_ID} -e AWS_SECRET_ACCESS_KEY=${AWS_SECRET_ACCESS_KEY} $(HELM_IMAGE)

GOOS                    := $(shell go env GOOS)
PSQL                    := docker exec -i test-server-db psql -U test_user -d test -v ON_ERROR_STOP=1

default: build

.PHONY: vendor
vendor:
	@GO111MODULE=on GOPROXY=direct GOSUMDB=off go mod tidy
	@GO111MODULE=on GOPROXY=direct GOSUMDB=off go mod vendor

.PHONY: gen
gen:
	@echo generate files...
	@$(GENTOOL) $(REPO)/pkg/server/pb/service.proto

	@$(GORUN) cmd/apidoc/utils/filter/filter.go --tags=$(SWAGGER_EXCLUDE_LIST) \
		--scheme=$(SWAGGER_FILE) > $(SWAGGER_FILE).patched

	@$(GORUN) cmd/apidoc/utils/adjust_schema.go --scheme=$(SWAGGER_FILE).patched \
		--apiKey=true --proto=http > pkg/server/pb/service.local.swagger.json

	@$(GORUN) cmd/apidoc/utils/adjust_schema.go --scheme=$(SWAGGER_FILE).patched \
		--apiKey=false --proto=https > $(SWAGGER_FILE)

	@rm $(SWAGGER_FILE).patched

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
	@mkdir -p $(BUILD_PATH)
	$(GOBUILD) -o $(SERVER_BINARY) $(SERVER_PATH)

.PHONY: test
test: fmt
	@go test $(GO_TEST_FLAGS) -tags=unit $(SERVER_PATH)

.PHONY: build
build: fmt build-server

.PHONY: image
image: server-image db-migration-image image-test-db

.PHONY: server-image
server-image:
	@docker build -f $(SERVER_DOCKERFILE) -t $(SERVER_IMAGE_FULL):$(IMAGE_VERSION) .
	@docker tag $(SERVER_IMAGE_FULL):$(IMAGE_VERSION) $(SERVER_IMAGE_FULL):latest

.PHONY: db-migration-image
db-migration-image:
	@docker build -f $(DB_MIGRATION_DOCKERFILE) -t $(DB_MIGRATION_IMAGE_FULL):$(IMAGE_VERSION) .
	@docker tag $(DB_MIGRATION_IMAGE_FULL):$(IMAGE_VERSION) $(DB_MIGRATION_IMAGE_FULL):latest

.PHONY: image-test-db
image-test-db:
	@docker build -f $(CUR_DIR)/db/Dockerfile.tests -t test-server-db:latest .

.PHONY: up
up:
	@docker-compose up -d

.PHONY: down
down:
	@docker-compose down -v

.PHONY: init-test-db
init-test-db:
	@sh $(CUR_DIR)/tests/wait_db.sh
	@cat $(CUR_DIR)/tests/testdb/cleanup_test_db.sql | $(PSQL) > /dev/null 2>&1
	@cat $(CUR_DIR)/tests/testdb/initialize_test_db.sql | $(PSQL) > /dev/null 2>&1

# test env
.PHONY: test-env
test-env: down build-server server-image image-test-db up init-test-db

.PHONY: push
push:
	@docker push $(SERVER_IMAGE_FULL):$(IMAGE_VERSION)
	@docker push $(DB_MIGRATION_IMAGE_FULL):$(IMAGE_VERSION)

.PHONY: helm-yaml
helm-yaml:
	@helm template $(CHART_PATH)/$(CHART) \
		--namespace $(NAMESPACE) \
		--name-template=$(NAMESPACE) \
		--set image.tagService=$(IMAGE_VERSION) \
		> $(BUILD_PATH)/$(CHART).yaml

.PHONY: push-chart
push-chart:
	$(HELM_GENERATOR) lint /repo/$(CHART)
	$(HELM_GENERATOR) package /repo/$(CHART) --version $(IMAGE_VERSION) -d /repo
	sed 's/{CHART_FILE}/$(CHART_FILE)/g' $(CHART_PATH)/build.properties.in > $(CHART_PATH)/build.properties
	$(HELM_GENERATOR) s3 push --force /repo/$(CHART_FILE) reponame

.PHONY: show-image-version
show-image-version:
	@echo $(IMAGE_VERSION)

.PHONY: clean
clean:
	$(BUILDER) sh -c "test ! -d $(SRCROOT)/bin || chmod -R 777 $(SRCROOT)/bin"
	@rm -rf "$(BIN_DIR)"
	@docker rmi -f $(shell docker images -q $(SERVER_IMAGE_FULL)) || true
	@docker rmi -f $(shell docker images -q $(DB_MIGRATION_IMAGE_FULL)) || true



