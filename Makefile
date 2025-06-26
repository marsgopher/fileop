# Ensure GOBIN is not set during build so that tools are installed to the correct path
unexport GOBIN

# GO related variables
GO           ?= go
GOFMT        ?= $(GO)fmt
GOTEST       := $(GO) test
FIRST_GOPATH := $(firstword $(subst :, ,$(shell $(GO) env GOPATH)))
GOOPTS       ?=
GOHOSTOS     ?= $(shell $(GO) env GOHOSTOS)
GOHOSTARCH   ?= $(shell $(GO) env GOHOSTARCH)
GO111MODULE  ?= $(shell $(GO) env GO111MODULE)
PKGS         := ./...
BUILD_ENV    ?=
BUILD_OPTS   ?= -trimpath


# Build arguments
APP          ?= $(shell basename "$(PWD)")
TAG          ?= unknown
APP_NAME     := $(APP)-$(TAG)
TIMESTAMP    := $(shell date +%Y-%m-%dT%T%z)
TESTOUT      := $(PWD)/coverage.out
BIN          ?= $(PWD)/bin
MAIN_PATH    ?= $(PWD)/.

# Get git reversion
REVER := unknown
ifneq ($(wildcard .git),)
REVER := $(shell git log -1 --pretty=format:"%H")
endif

ifeq ($(GOHOSTARCH),amd64 arm64)
	ifeq ($(GOHOSTOS),$(filter $(GOHOSTOS),linux freebsd darwin windows))
		# Only supported on amd64/arm64
		test-flags := -race
	endif
endif

GOLANGCI_LINT :=
GOLANGCI_LINT_OPTS ?=
GOLANGCI_LINT_VERSION ?= v1.64.8
# golangci-lint only supports linux, darwin and windows platforms on i386/amd64/arm64.
# windows isn't included here because of the path separator being different.
ifeq ($(GOHOSTOS),$(filter $(GOHOSTOS),linux darwin))
	ifeq ($(GOHOSTARCH),$(filter $(GOHOSTARCH),amd64 i386 arm64))
		GOLANGCI_LINT := $(FIRST_GOPATH)/bin/golangci-lint
	endif
endif

.PHONY: all
all: test

## test: running tests
.PHONY: test
test: lint
	@echo ">> running tests"
	$(BUILD_ENV) $(GOTEST) $(GOOPTS) $(test-flags) -cover $(PKGS) -coverprofile $(TESTOUT)

## lint: running code inspection
.PHONY: lint
lint: $(GOLANGCI_LINT)
ifdef GOLANGCI_LINT
	@echo ">> running golangci-lint"
ifdef GO111MODULE
# 'go list' needs to be executed before staticcheck to prepopulate the modules cache.
# Otherwise staticcheck might fail randomly for some reason not yet explained.
	GO111MODULE=$(GO111MODULE) $(GO) list -e -compiled -test=true -export=false -deps=true -find=false -tags= -- ./... > /dev/null
	GO111MODULE=$(GO111MODULE) $(GOLANGCI_LINT) run $(GOLANGCI_LINT_OPTS) $(PKGS)
else
	$(GOLANGCI_LINT) run $(GOLANGCI_LINT_OPTS) $(PKGS)
endif
endif

# install golangci-lint if not exist
ifdef GOLANGCI_LINT
$(GOLANGCI_LINT):
	mkdir -p $(FIRST_GOPATH)/bin
	curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/$(GOLANGCI_LINT_VERSION)/install.sh \
		| sed -e '/install -d/d' \
		| sh -s -- -b $(FIRST_GOPATH)/bin $(GOLANGCI_LINT_VERSION)
endif
