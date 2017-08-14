COVERAGEDIR = coverage
SERVICE=local
ifdef CIRCLE_WORKING_DIRECTORY
  COVERAGEDIR = $(CIRCLE_WORKING_DIRECTORY)/coverage
	SERVICE=circle-ci
endif

PACKAGES='./generator' './example'

.PHONY: all
all: generate fmt build test example cover install 

.PHONY: install-deps
install-deps:
	go get github.com/golang/dep/cmd/dep
	go get -v github.com/jteeuwen/go-bindata/...
	go get -v golang.org/x/tools/cmd/cover
	go get -v github.com/mattn/goveralls
	go get -v github.com/modocache/gover
	dep ensure

build:
	go generate ./generator
	if [ ! -d bin ]; then mkdir bin; fi
	go build -v -o bin/go-enum .

fmt:
	gofmt -l -w -s $$(find . -type f -name '*.go' -not -path "./vendor/*")

test: generate gen-test
	if [ ! -d coverage ]; then mkdir coverage; fi
	go test -v ./generator -race -cover -coverprofile=$(COVERAGEDIR)/generator.coverprofile

cover: gen-test test
	go tool cover -html=$(COVERAGEDIR)/generator.coverprofile -o $(COVERAGEDIR)/generator.html

tc: test cover
coveralls:
	if [ ! -d $(COVERAGEDIR) ]; then mkdir $(COVERAGEDIR); fi
	gover $(COVERAGEDIR) $(COVERAGEDIR)/coveralls.coverprofile
	goveralls -coverprofile=$(COVERAGEDIR)/coveralls.coverprofile -service=$(SERVICE) -repotoken=$(COVERALLS_TOKEN)

clean:
	go clean
	rm -f bin/go-enum
	rm -rf coverage/

.PHONY: generate
generate:
	go generate $(PACKAGES)

gen-test: build install
	go generate $(PACKAGES)

install:
	go install

phony: clean tc build

.PHONY: example
example:
	go generate ./example