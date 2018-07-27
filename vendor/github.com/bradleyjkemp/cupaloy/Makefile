.PHONY: install
install: get_dependencies install_linters

.PHONY: get_dependencies
get_dependencies:
	go get github.com/mattn/goveralls
	go get github.com/golang/dep/cmd/dep
	$(GOPATH)/bin/dep ensure

.PHONY: install_linters
install_linters:
	go get github.com/alecthomas/gometalinter
	$(GOPATH)/bin/gometalinter --install

.PHONY: lint
lint:
	$(GOPATH)/bin/gometalinter --vendor ./...

.PHONY: test
test: lint
	go test ./...


.PHONY: test-ci
test-ci: coverage lint

.PHONY: coverage
coverage:
	go test -v -coverpkg ./... -coverprofile coverage.out ./...

.PHONY: clean
clean:
	rm -rf examples/ignored*
