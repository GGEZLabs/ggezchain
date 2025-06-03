#!/usr/bin/make -f

include contrib/devtools/Makefile
include contrib/devtools/lint.mk

VERSION := $(shell echo $(shell git describe --tags) | sed 's/^v//')
COMMIT := $(shell git log -1 --format='%H')
LEDGER_ENABLED ?= true
# for dockerized protobuf tools
DOCKER := $(shell which docker)
OUTPUT_DIR := build
PLATFORMS = linux/amd64 darwin/amd64 darwin/arm64
APPNAME := ggezchain

export GO111MODULE = on

# don't override user values
ifeq (,$(VERSION))
  VERSION := $(shell git describe --tags --always)
  # if VERSION is empty, then populate it with branch's name and raw commit hash
  ifeq (,$(VERSION))
    VERSION := $(BRANCH)-$(COMMIT)
  endif
endif

# process build tags

build_tags = netgo
ifeq ($(LEDGER_ENABLED),true)
  ifeq ($(OS),Windows_NT)
    GCCEXE = $(shell where gcc.exe 2> NUL)
    ifeq ($(GCCEXE),)
      $(error gcc.exe not installed for ledger support, please install or set LEDGER_ENABLED=false)
    else
      build_tags += ledger
    endif
  else
    UNAME_S = $(shell uname -s)
    ifeq ($(UNAME_S),OpenBSD)
      $(warning OpenBSD detected, disabling ledger support (https://github.com/cosmos/cosmos-sdk/issues/1988))
    else
      GCC = $(shell command -v gcc 2> /dev/null)
      ifeq ($(GCC),)
        $(error gcc not installed for ledger support, please install or set LEDGER_ENABLED=false)
      else
        build_tags += ledger
      endif
    endif
  endif
endif

ifeq ($(WITH_CLEVELDB),yes)
  build_tags += gcc
endif
build_tags += $(BUILD_TAGS)
build_tags := $(strip $(build_tags))

whitespace :=
empty = $(whitespace) $(whitespace)
comma := ,
build_tags_comma_sep := $(subst $(empty),$(comma),$(build_tags))

# flags '-s -w' resolves an issue with xcode 16 and signing of go binaries
# ref: https://github.com/golang/go/issues/63997
ldflags = -X github.com/cosmos/cosmos-sdk/version.Name=ggezchain \
		  -X github.com/cosmos/cosmos-sdk/version.AppName=ggezchaind \
		  -X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
		  -X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT) \
		  -X "github.com/cosmos/cosmos-sdk/version.BuildTags=$(build_tags_comma_sep)" \
		  -s -w

ifeq ($(WITH_CLEVELDB),yes)
  ldflags += -X github.com/cosmos/cosmos-sdk/types.DBBackend=cleveldb
endif
ifeq ($(LINK_STATICALLY),true)
	ldflags += -linkmode=external -extldflags "-Wl,-z,muldefs -static"
endif
ldflags += $(LDFLAGS)
ldflags := $(strip $(ldflags))

BUILD_FLAGS := -tags "$(build_tags_comma_sep)" -ldflags '$(ldflags)' -trimpath

all: install lint test

build: go.sum
ifeq ($(OS),Windows_NT)
	$(error wasmd server not supported. Use "make build-windows-client" for client)
	exit 1
else
	go build -mod=readonly $(BUILD_FLAGS) -o $(OUTPUT_DIR)/ggezchaind ./cmd/ggezchaind
endif

build-windows-client: go.sum
	GOOS=windows GOARCH=amd64 go build -mod=readonly $(BUILD_FLAGS) -o $(OUTPUT_DIR)/ggezchaind.exe ./cmd/ggezchaind

install: go.sum
	go install -mod=readonly $(BUILD_FLAGS) ./cmd/ggezchaind

build-image:
	docker build -f Dockerfile -t ggezlabs/ggezchain .

mocks:
	@go install go.uber.org/mock/mockgen@v0.5.0
	sh ./scripts/mockgen.sh
.PHONY: mocks

########################################
### Tools & dependencies
########################################

go-mod-cache: go.sum
	@echo "--> Download go modules to local cache"
	@go mod download

go.sum: go.mod
	@echo "--> Ensure dependencies have not been modified"
	@go mod verify

clean:
	rm -rf $(OUTPUT_DIR)/*

.PHONY: all install \
	go-mod-cache clean build \
    build-windows-client

########################################
### Testing
########################################
PACKAGES_E2E=$(shell cd tests/e2e && go list ./... | grep '/e2e')
PACKAGES_UNIT=$(shell go list ./... | grep -v -e '/tests/e2e')

test-all: test test-race test-cover

test:
	@VERSION=$(VERSION) go test -mod=readonly -tags='ledger test_ledger_mock'  $(PACKAGES_UNIT)

test-race:
	@VERSION=$(VERSION) go test -mod=readonly -race -tags='ledger test_ledger_mock'  $(PACKAGES_UNIT)

test-cover:
	@go test -mod=readonly -timeout 30m -race -coverprofile=coverage.txt -covermode=atomic -tags='ledger test_ledger_mock' ./...

test-e2e: build-image
	@VERSION=$(VERSION) go test -mod=readonly -timeout=35m -v $(PACKAGES_E2E)

benchmark:
	@go test -mod=readonly -bench=. $(PACKAGES_UNIT)

.PHONY: test test-all \
	test test-race \
	test-cover benchmark

###############################################################################
###                                Linting                                  ###
###############################################################################

# golangci_lint_cmd=golangci-lint
# golangci_version=v1.61.0

# lint:
# 	@echo "--> Running linter"
# 	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(golangci_version)
# 	@$(golangci_lint_cmd) run ./... --timeout 15m

format-tools:
	go install mvdan.cc/gofumpt@v0.4.0
	go install github.com/client9/misspell/cmd/misspell@v0.3.4
	go install github.com/daixiang0/gci@v0.11.2

format: format-tools
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "./tests/system/vendor*" -not -path "*.git*" -not -path "./client/lcd/statik/statik.go" | xargs gofumpt -w
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "./tests/system/vendor*" -not -path "*.git*" -not -path "./client/lcd/statik/statik.go" | xargs misspell -w
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "./tests/system/vendor*" -not -path "*.git*" -not -path "./client/lcd/statik/statik.go" | xargs gci write --skip-generated -s standard -s default -s "prefix(cosmossdk.io)" -s "prefix(github.com/cosmos/cosmos-sdk)" -s "prefix(github.com/CosmWasm/wasmd)" --custom-order

mod-tidy:
	go mod tidy

.PHONY: format-tools lint format mod-tidy


###############################################################################
###                                Protobuf                                 ###
###############################################################################
CURRENT_UID := $(shell id -u)
CURRENT_GID := $(shell id -g)

protoVer=0.13.2
protoImageName=ghcr.io/cosmos/proto-builder:$(protoVer)
protoImage=sudo "$(DOCKER)" run -e BUF_CACHE_DIR=/tmp/buf --rm -v "$(CURDIR)":/workspace:rw --user ${CURRENT_UID}:${CURRENT_GID} --workdir /workspace $(protoImageName)

proto-gen:
	@echo "Generating protobuf files..."
	@ignite generate proto-go --yes

proto-format:
	@echo "Formatting Protobuf files"
	@$(protoImage) find ./ -name "*.proto" -exec clang-format -i {} \;

.PHONY: proto-gen proto-format

###############################################################################
###                                    testnet                              ###
###############################################################################

setup-testnet: mod-tidy set-testnet-configs setup-testnet-keys

# Run this before testnet keys are added
# This chain id is used in the testnet.json as well
set-testnet-configs:
	ggezchaind config set client chain-id localchain_9000-1
	ggezchaind config set client keyring-backend test
	ggezchaind config set client output text

# import keys from testnet.json into test keyring
setup-testnet-keys:
	@echo "Adding acc0..."
	@echo "decorate bright ozone fork gallery riot bus exhaust worth way bone indoor calm squirrel merry zero scheme cotton until shop any excess stage laundry" | ggezchaind keys add acc0 --recover
	@echo "Adding acc1..."
	@echo "wealth flavor believe regret funny network recall kiss grape useless pepper cram hint member few certain unveil rather brick bargain curious require crowd raise" | ggezchaind keys add acc1 --recover
	
.PHONY: setup-testnet set-testnet-configs setup-testnet-keys