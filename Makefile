.PHONY: all
all: lint test

.PHONY: lint
lint: tools
	golangci-lint run

.PHONY: test
test: tools
	ginkgo -r -p -keepGoing -randomizeAllSpecs -randomizeSuites -trace -race -progress

sources = $(shell find . -name '*.go' -not -path './vendor/*')
.PHONY: goimports
goimports: tools
	$(GOBIN)/goimports -w $(sources)

export GOBIN ?= $(PWD)/bin
export PATH := $(GOBIN):$(PATH)
GOIMPORTS := $(GOBIN)/goimports
GOCOV := $(GOBIN)/gocov
GOCOV_HTML := $(GOBIN)/gocov-html
GOLANGCI_LINT := $(GOBIN)/golangci-lint
GINKGO := $(GOBIN)/ginkgo

$(GOIMPORTS):
	go install golang.org/x/tools/cmd/goimports

$(GOCOV):
	go install github.com/axw/gocov/gocov

$(GOCOV_HTML):
	go install github.com/matm/gocov-html

$(GOLANGCI_LINT):
	go install github.com/golangci/golangci-lint/cmd/golangci-lint

$(GINKGO):
	go install github.com/onsi/ginkgo/ginkgo

tools: $(GOIMPORTS) $(GINKGO) $(GOCOV) $(GOCOV_HTML) $(GOLANGCI_LINT)
