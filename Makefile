# Make dependency graph:
# test -> lint -> $(GOMETALINTER)

# Skip vendor directory.
PKGS := $(shell go list ./... | grep -v /vendor)
.PHONY: test
# lint before we test, `lint` is a pre-requisite of `test`.
test: lint
	go test $(PKGS)

BIN_DIR := $(GOPATH)/bin
GOMETALINTER := $(BIN_DIR)/gometalinter
.PHONY: lint
lint: $(GOMETALINTER)
	gometalinter ./... --vendor

# GOMETALINTER will be rebuilt if and only if the `gometalinter` binary is not present at BIN_DIR.
# Otherwise, it's up-to-date.
$(GOMETALINTER):
	go get -u github.com/alecthomas/gometalinter
	gometalinter --install &> /dev/null
