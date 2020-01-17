.DEFAULT_GOAL:=all
ifdef VERBOSE
V = -v
else
.SILENT:
endif

include $(wildcard *.mk)

COVERAGEDIR = coverage
SERVICE=local
ifdef CIRCLE_WORKING_DIRECTORY
  COVERAGEDIR = $(CIRCLE_WORKING_DIRECTORY)/coverage
	SERVICE=circle-ci
endif

PACKAGES='./generator' './example'

.PHONY: all
all: build fmt test example cover install

build: deps
	go generate ./generator
	if [ ! -d bin ]; then mkdir bin; fi
	go build -v -o bin/go-enum .

fmt:
	gofmt -l -w -s $$(find . -type f -name '*.go' -not -path "./vendor/*")

test: gen-test generate
	go test -v -race -coverprofile=coverage.out ./...

cover: gen-test test
	go tool cover -html=coverage.out -o coverage.html

tc: test cover
coveralls:
	goveralls -coverprofile=coverage.out -service=$(SERVICE) -repotoken=$(COVERALLS_TOKEN)

clean: cleandeps
	go clean
	rm -f bin/go-enum
	rm -rf coverage/

.PHONY: generate
generate:
	go generate $(PACKAGES)

gen-test: build
	$(GO) generate $(PACKAGES)

install:
	go install

phony: clean tc build

.PHONY: example
example:
	go generate ./example
