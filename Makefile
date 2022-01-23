.DEFAULT_GOAL:=all
ifdef VERBOSE
V = -v
else
.SILENT:
endif

GO ?= CGO_ENABLED=0 GO111MODULE=on go
COVERAGEDIR= coverage
SERVICE=local

ifdef GITHUB_ACTIONS
SERVICE=github-actions
endif

DATE := $(shell date -u '+%FT%T%z')
GITHUB_SHA ?= $(shell git rev-parse HEAD)
GITHUB_REF ?= local

LDFLAGS += -X "main.version=$(GITHUB_REF)"
LDFLAGS += -X "main.commit=$(GITHUB_SHA)"
LDFLAGS += -X "main.date=$(DATE)"
LDFLAGS += -X "main.builtBy=$(USER)"
LDFLAGS += -extldflags '-static'

define goinstall
	mkdir -p $(shell pwd)/bin
	echo "Installing $(1)"
	GOBIN=$(shell pwd)/bin go install $(1)
endef

GOBINDATA=bin/go-bindata
GOIMPORTS=bin/goimports
GOVERALLS=bin/goveralls
MOCKGEN=bin/mockgen
deps: $(MOCKGEN)
deps: $(GOBINDATA)
deps: $(GOIMPORTS)

PACKAGES='./generator' './example'

.PHONY: all
all: build fmt test example cover install

build: deps
	$(GO) generate ./generator
	if [ ! -d bin ]; then mkdir bin; fi
	$(GO) build -v -o bin/go-enum -ldflags='-X "main.version=example" -X "main.commit=example" -X "main.date=example" -X "main.builtBy=example"'  .

fmt:
	gofmt -l -w -s $$(find . -type f -name '*.go' -not -path "./vendor/*")

test: gen-test generate
	$(GO) test -v -coverprofile=coverage.out ./...

cover: gen-test test
	$(GO) tool cover -html=coverage.out -o coverage.html

tc: test cover
coveralls: $(GOVERALLS)
	$(GOVERALLS) -coverprofile=coverage.out -service=$(SERVICE) -repotoken=$(COVERALLS_TOKEN)

clean:
	$(GO) clean
	rm -f bin/go-enum
	rm -rf coverage/
	rm -rf bin/
	rm -rf dist/

.PHONY: generate
generate:
	$(GO) generate $(PACKAGES)

gen-test: build
	$(GO) generate $(PACKAGES)

install:
	$(GO) install

phony: clean tc build

.PHONY: example
example:
	$(GO) generate ./example/...

bin/goimports: go.sum
	$(call goinstall,golang.org/x/tools/cmd/goimports)

bin/mockgen: go.sum
	$(call goinstall,github.com/golang/mock/mockgen)

bin/goveralls: go.sum
	$(call goinstall,github.com/mattn/goveralls)

bin/go-bindata: go.sum
	$(call goinstall,github.com/kevinburke/go-bindata/go-bindata)
