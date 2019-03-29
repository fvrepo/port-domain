BASEPATH = $(shell pwd)
export PATH := $(BASEPATH)/bin:$(PATH)

# Basic go commands
GOCMD      = go
GOBUILD    = $(GOCMD) build
GOINSTALL  = $(GOCMD) install
GORUN      = $(GOCMD) run
GOCLEAN    = $(GOCMD) clean
GOTEST     = $(GOCMD) test
GOGET      = $(GOCMD) get
GOFMT      = $(GOCMD) fmt
GOGENERATE = $(GOCMD) generate
GOTYPE     = $(GOCMD)type

# Docker
DOCKER_COMPOSE = docker-compose

# GRPC
PROTOC       = protoc
PROTOCGOGEN  = protoc-gen-go

BINARY = port

BUILD_DIR = $(BASEPATH)
COVERAGE_DIR  = $(BUILD_DIR)/coverage
SUBCOV_DIR    = $(COVERAGE_DIR)/packages

# all src packages without vendor and generated code
PKGS = $(shell go list ./... | grep -v /vendor | grep -v /internal/server/grpcapi)

# Colors
GREEN_COLOR   = \033[0;32m
PURPLE_COLOR  = \033[0;35m
DEFAULT_COLOR = \033[m

all: clean fmt build test

help:
	@echo 'Usage: make <TARGETS> ... <OPTIONS>'
	@echo ''
	@echo 'Available targets are:'
	@echo ''
	@echo '    help               Show this help screen.'
	@echo '    clean              Remove binary.'
	@echo '    test               Run unit tests.'
	@echo '    lint               Run all linters including vet and gosec and others'
	@echo '    coverage           Report code tests coverage.'
	@echo '    fmt                Run gofmt on package sources.'
	@echo '    build              Compile packages and dependencies.'
	@echo '    grpc               Generate pb.go from proto file'
	@echo '    version            Print Go version.'
	@echo ''
	@echo 'Targets run by default are: clean fmt lint test.'
	@echo ''

clean:
	@echo " $(GREEN_COLOR)[clean]$(DEFAULT_COLOR)"
	@$(GOCLEAN)
	@if [ -f $(BINARY) ] ; then rm $(BINARY) ; fi

coverage:
	@echo " [$(GREEN_COLOR)coverage$(DEFAULT_COLOR)]"
	@-mkdir -p $(SUBCOV_DIR)/
	@for package in $(PKGS); do $(GOTEST) -coverpkg=./... -coverprofile $(SUBCOV_DIR)/`basename "$$package"`.cov "$$package"; done
	@echo 'mode: count' > $(COVERAGE_DIR)/coverage.cov ;
	@tail -q -n +2 $(SUBCOV_DIR)/*.cov >> $(COVERAGE_DIR)/coverage.cov ;
	@go tool cover -func=$(COVERAGE_DIR)/coverage.cov ;
	@if [ $(html) ]; then go tool cover -html=$(COVERAGE_DIR)/coverage.cov -o coverage.html ; fi
	@rm -rf $(COVERAGE_DIR);

lint:
	@echo " [$(GREEN_COLOR)lint$(DEFAULT_COLOR)]"
	@$(GORUN) ./vendor/github.com/golangci/golangci-lint/cmd/golangci-lint/main.go run \
	--no-config --enable=errcheck --enable=gosec --enable=gocyclo --enable=nakedret \
	--enable=prealloc --enable=gofmt --enable=goimports --enable=megacheck --enable=misspell ./...

test:
	@echo " $(GREEN_COLOR)[test]$(DEFAULT_COLOR)"
	@$(GOTEST) -race $(PKGS)

fmt:
	@echo " $(GREEN_COLOR)[format]$(DEFAULT_COLOR)"
	@$(GOFMT) $(PKGS)

build:
	@echo " $(GREEN_COLOR)[build]$(DEFAULT_COLOR)"
	$(GOBUILD) --tags static -o $(BINARY)

version:
	@echo " $(GREEN_COLOR)[version]$(DEFAULT_COLOR)"
	@$(GOCMD) version

grpc:
	@mkdir -p ./bin
ifeq ("$(wildcard ./bin/$(PROTOCGOGEN))","")
	@echo " $(PURPLE_COLOR)[build protoc-go-gen]$(DEFAULT_COLOR)"
	@$(GOBUILD) -o ./bin/$(PROTOCGOGEN) ./vendor/github.com/golang/protobuf/protoc-gen-go/
endif

	@echo " [$(GREEN_COLOR)grpc$(DEFAULT_COLOR)]"
		@-rm -rf ./pkg/grpcapi
		@mkdir -p ./pkg/grpcapi/port

	@${PROTOC} \
		-I ./api \
		./api/port-grpc.proto \
		--go_out=plugins=grpc:./pkg/grpcapi/port