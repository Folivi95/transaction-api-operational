include .bingo/Variables.mk

.PHONY: test
default: test

setup:
	git config core.hooksPath .hooks
	git config --global url."git@github.com:saltpay".insteadOf https://github.com/saltpay
	go install github.com/bwplotka/bingo@latest
	brew install -q ctlptl helm kind kubernetes-cli tilt yq

dev:
	ctlptl apply -f ./local/kind-cluster.yaml
	tilt up

stop:
	@tilt down

# use this when tilt is very unhappy and you have no idea why
kill:
	make stop
	# this script takes a while, so be patient...
	ctlptl delete cluster kind-kind
	docker ps -aq | xargs docker rm -f
	ctlptl docker-desktop quit
	ctlptl docker-desktop open
	@echo '💀  Killed everything. Wait for docker to come back up before running `make dev`'

enable-o11y:
	tilt args -- --o11y

t: test
test: lint unit-tests integration-tests acceptance-tests

unit-tests:
	$(GOTEST) -shuffle=on --tags=unit ./...

integration-tests:
	$(GOTEST) -count=1 --tags=integration ./...

acceptance-tests:
	$(GOTEST) -count=1 --tags=acceptance ./...

race-condition-tests:
	$(GOTEST) -count=1 --tags=race ./...

lint:
	$(GOLANGCI_LINT) run --timeout=5m

lf: lint_fix
lint_fix:
	@$(GOLANGCI_LINT) run ./... --fix

generate:
	@go generate ./...

bump-chart-version:
	$(CHARTVERSION) bump --remote=origin --trunk=main

# this is used by Tilt. It's in the makefile to make use of Make's caching/smart capabilities
local/Generated_DoNotModify.Local.Dockerfile: Dockerfile
	cat Dockerfile | sed 's/USER nobody//' > local/Generated_DoNotModify.Local.Dockerfile

mod:
	go mod vendor -v

tidy:
	go mod tidy -v

clean-n:
	go clean -n

clean-x:
	go clean -x

clean-modcache:
	go clean -modcache
