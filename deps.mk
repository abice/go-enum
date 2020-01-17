GO ?= go

GOBINDATA=bin/go/github.com/kevinburke/go-bindata/go-bindata
GOIMPORTS=bin/go/golang.org/x/tools/cmd/goimports
MOCKGEN=bin/go/github.com/golang/mock/mockgen
TOOLS = $(GOBINDATA) \
	$(GOIMPORTS) \
	$(MOCKGEN)

cleandeps:
	if [ -d "./bin" ]; then rm -rf "./bin"; fi

freshdeps: cleandeps deps

deps: $(TOOLS)
bin/go/%:
	@echo "installing $*"
	$(GO) build -o bin/go/$* $*
