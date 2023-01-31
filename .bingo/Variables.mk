# Auto generated binary variables helper managed by https://github.com/bwplotka/bingo v0.6. DO NOT EDIT.
# All tools are designed to be build inside $GOBIN.
BINGO_DIR := $(dir $(lastword $(MAKEFILE_LIST)))
GOPATH ?= $(shell go env GOPATH)
GOBIN  ?= $(firstword $(subst :, ,${GOPATH}))/bin
GO     ?= $(shell which go)

# Below generated variables ensure that every time a tool under each variable is invoked, the correct version
# will be used; reinstalling only if needed.
# For example for bingo variable:
#
# In your main Makefile (for non array binaries):
#
#include .bingo/Variables.mk # Assuming -dir was set to .bingo .
#
#command: $(BINGO)
#	@echo "Running bingo"
#	@$(BINGO) <flags/args..>
#
BINGO := $(GOBIN)/bingo-v0.6.0
$(BINGO): $(BINGO_DIR)/bingo.mod
	@# Install binary/ries using Go 1.14+ build command. This is using bwplotka/bingo-controlled, separate go module with pinned dependencies.
	@echo "(re)installing $(GOBIN)/bingo-v0.6.0"
	@cd $(BINGO_DIR) && $(GO) build -mod=mod -modfile=bingo.mod -o=$(GOBIN)/bingo-v0.6.0 "github.com/bwplotka/bingo"

CHARTVERSION := $(GOBIN)/chartversion-v0.2.3
$(CHARTVERSION): $(BINGO_DIR)/chartversion.mod
	@# Install binary/ries using Go 1.14+ build command. This is using bwplotka/bingo-controlled, separate go module with pinned dependencies.
	@echo "(re)installing $(GOBIN)/chartversion-v0.2.3"
	@cd $(BINGO_DIR) && $(GO) build -mod=mod -modfile=chartversion.mod -o=$(GOBIN)/chartversion-v0.2.3 "github.com/saltpay/go-utils/cmd/chartversion"

GOLANGCI_LINT := $(GOBIN)/golangci-lint-v1.47.3
$(GOLANGCI_LINT): $(BINGO_DIR)/golangci-lint.mod
	@# Install binary/ries using Go 1.14+ build command. This is using bwplotka/bingo-controlled, separate go module with pinned dependencies.
	@echo "(re)installing $(GOBIN)/golangci-lint-v1.47.3"
	@cd $(BINGO_DIR) && $(GO) build -mod=mod -modfile=golangci-lint.mod -o=$(GOBIN)/golangci-lint-v1.47.3 "github.com/golangci/golangci-lint/cmd/golangci-lint"

GOTEST := $(GOBIN)/gotest-v0.0.6
$(GOTEST): $(BINGO_DIR)/gotest.mod
	@# Install binary/ries using Go 1.14+ build command. This is using bwplotka/bingo-controlled, separate go module with pinned dependencies.
	@echo "(re)installing $(GOBIN)/gotest-v0.0.6"
	@cd $(BINGO_DIR) && $(GO) build -mod=mod -modfile=gotest.mod -o=$(GOBIN)/gotest-v0.0.6 "github.com/rakyll/gotest"

MOQ := $(GOBIN)/moq-v0.2.7
$(MOQ): $(BINGO_DIR)/moq.mod
	@# Install binary/ries using Go 1.14+ build command. This is using bwplotka/bingo-controlled, separate go module with pinned dependencies.
	@echo "(re)installing $(GOBIN)/moq-v0.2.7"
	@cd $(BINGO_DIR) && $(GO) build -mod=mod -modfile=moq.mod -o=$(GOBIN)/moq-v0.2.7 "github.com/matryer/moq"

SELECT := $(GOBIN)/select-v0.2.0
$(SELECT): $(BINGO_DIR)/select.mod
	@# Install binary/ries using Go 1.14+ build command. This is using bwplotka/bingo-controlled, separate go module with pinned dependencies.
	@echo "(re)installing $(GOBIN)/select-v0.2.0"
	@cd $(BINGO_DIR) && $(GO) build -mod=mod -modfile=select.mod -o=$(GOBIN)/select-v0.2.0 "github.com/tamj0rd2/coauthor-select/cmd/select"

VALIDATE := $(GOBIN)/validate-v0.2.0
$(VALIDATE): $(BINGO_DIR)/validate.mod
	@# Install binary/ries using Go 1.14+ build command. This is using bwplotka/bingo-controlled, separate go module with pinned dependencies.
	@echo "(re)installing $(GOBIN)/validate-v0.2.0"
	@cd $(BINGO_DIR) && $(GO) build -mod=mod -modfile=validate.mod -o=$(GOBIN)/validate-v0.2.0 "github.com/tamj0rd2/coauthor-select/cmd/validate"

