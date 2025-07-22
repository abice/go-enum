.DEFAULT_GOAL:=all
ifdef VERBOSE
V = -v
else
.SILENT:
endif

GO ?= go
COVERAGEDIR= coverage
SERVICE=local

ifdef GITHUB_ACTIONS
SERVICE=github
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

GOIMPORTS=bin/goimports
GOVERALLS=bin/goveralls
MOCKGEN=bin/mockgen
deps: $(MOCKGEN)
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
	$(GO) test -v -race -shuffle on -coverprofile=coverage.out ./...
	$(GO) test -v -race -shuffle on --tags=example ./example

cover: gen-test test
	$(GO) tool cover -html=coverage.out -o coverage.html

tc: test cover
coveralls: $(GOVERALLS)
	$(GOVERALLS) -coverprofile=coverage.out -service=$(SERVICE) -repotoken=$(COVERALLS_TOKEN)

clean:
	rm -f bin/go-enum
	rm -rf coverage/
	rm -rf bin/
	rm -rf dist/
	$(GO) clean

.PHONY: assert-no-changes
assert-no-changes:
	@if [ -n "$(shell git status --porcelain)" ]; then \
		echo "git changes found: $(shell git status --porcelain)"; \
		exit 1; \
	fi

.PHONY: generate
generate:
	$(GO) generate --tags=example $(PACKAGES)

gen-test: build
	$(GO) generate --tags=example $(PACKAGES)

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

# snapshots: snapshots_1.17
snapshots: snapshots_1.24

snapshots_%: clean
	echo "##### updating snapshots for golang $* #####"
	docker run -i -t -w /app -v $(shell pwd):/app --entrypoint /bin/sh golang:$* -c './update-snapshots.sh || true'

.PHONY: ci
ci: docker_1.23
ci: docker_1.24

docker_%:
	echo "##### testing golang $* #####"
	docker run -i -t -w /app -v $(shell pwd):/app --entrypoint /bin/sh golang:$* -c 'make clean && make'

.PHONY: pullimages
pullimages: pullimage_1.23
pullimages: pullimage_1.24

pullimage_%:
	docker pull golang:$*

build_docker:
	KO_DOCKER_REPO=abice/go-enum VERSION=$(GITHUB_REF) COMMIT=$(GITHUB_SHA) DATE=$(DATE) BUILT_BY=$(USER) ko build --bare --local
