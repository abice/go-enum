GOBINDATA=bin/go/github.com/kevinburke/go-bindata/go-bindata
GOIMPORTS=bin/go/golang.org/x/tools/cmd/goimports
GOVERALLS=bin/go/github.com/mattn/goveralls
MOCKGEN=bin/go/github.com/golang/mock/mockgen
TOOLS = $(GOBINDATA) \
	$(GOIMPORTS) \
	$(MOCKGEN) \
	$(GOVERALLS)

cleandeps:
	if [ -d "./bin" ]; then rm -rf "./bin"; fi

freshdeps: cleandeps deps

deps: $(TOOLS)
bin/go/%:
	@echo "installing $*"
	$(GO) build -o bin/go/$* $*
