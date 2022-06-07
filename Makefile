.DEFAULT_GOAL:=all
ifdef VERBOSE
V = -v
else
.SILENT:
endif

GO ?= GO111MODULE=on go
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
	-$(GO) fmt ./...

test: gen-test generate
	$(GO) test -v -race -coverprofile=coverage.out ./...

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

.PHONY: assert-no-changes
assert-no-changes:
	@if [ -n "$(shell git status --porcelain)" ]; then \
		echo "git changes found: $(shell git status --porcelain)"; \
		exit 1; \
	fi

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

generate1_15: generator/assets/assets.go generator/enum.tmpl
	docker run -i -t -w /app -v $(shell pwd):/app --entrypoint /bin/sh golang:1.15 -c 'make clean $(GOBINDATA) && $(GO) generate ./generator && make clean'

.PHONY: snapshots1_18
snapshots1_18:
	docker run -i -t -w /app -v $(shell pwd):/app --entrypoint /bin/sh golang:1.18 -c './update-snapshots.sh || true && make clean && make'

.PHONY: ci
ci: docker_1.14
ci: docker_1.15
ci: docker_1.16
ci: docker_1.17
ci: docker_1.18

docker_%:
	echo "##### testing golang $* #####"
	docker run -i -t -w /app -v $(shell pwd):/app --entrypoint /bin/sh golang:$* -c 'make clean && make'
