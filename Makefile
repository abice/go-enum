COVERAGEDIR = coverage
SERVICE=local
ifdef CIRCLE_ARTIFACTS
  COVERAGEDIR = $(CIRCLE_ARTIFACTS)
	SERVICE=circle-ci
endif

.PHONY: all
all: generate fmt build test example cover install 

.PHONY: install-deps
install-deps:
	glide install

build: generate
	if [ ! -d bin ]; then mkdir bin; fi
	go build -v -o bin/go-enum .

fmt:
	gofmt -l -w -s $$(find . -type f -name '*.go' -not -path "./vendor/*")

test: generate gen-test
	if [ ! -d coverage ]; then mkdir coverage; fi
	go test -v ./generator -race -cover -coverprofile=$(COVERAGEDIR)/generator.coverprofile

cover:
	go tool cover -html=$(COVERAGEDIR)/generator.coverprofile -o $(COVERAGEDIR)/generator.html

tc: test cover
coveralls:
	gover $(COVERAGEDIR) $(COVERAGEDIR)/coveralls.coverprofile
	goveralls -coverprofile=$(COVERAGEDIR)/coveralls.coverprofile -service=$(SERVICE) -repotoken=$(COVERALLS_TOKEN)
clean:
	go clean
	rm -f bin/go-enum
	rm -rf coverage/

.PHONY: generate
generate:
	go generate $$(glide nv)

gen-test: build install
	go generate $$(glide nv)

install:
	go install

phony: clean tc build

.PHONY: example
example:
	go generate ./example