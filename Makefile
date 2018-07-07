# Make dependency graph:
# test -> lint -> $(GOMETALINTER)

# PKGS stores the results of the script in $().
# i.e., the script results in: news/api, the rule test then will append the package name to
# `go test`, like so: `go test news/api`. (skipping /vendor directory)
PKGS := $(shell go list ./... | grep -v /vendor)
.PHONY: test
# lint before we test, `lint` is a pre-requisite of `test`.
test: lint
	go test $(PKGS)

.PHONY: lint
lint: $(GOMETALINTER)
	gometalinter ./... --vendor

BIN_DIR := $(GOPATH)/bin
GOMETALINTER := $(BIN_DIR)/gometalinter

$(GOMETALINTER):
	go get -u github.com/alecthomas/gometalinter
	gometalinter --install &> /dev/null
