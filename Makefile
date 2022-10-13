GOCMD:=$(shell which go)
GOLINT:=$(shell which golint)
GOIMPORT:=$(shell which goimports)
GOFMT:=$(shell which gofmt)
GOINSTALL:=$(GOCMD) install
GOCLEAN:=$(GOCMD) clean
GOTEST:=$(GOCMD) test
GOGET:=$(GOCMD) get
GOLIST:=$(GOCMD) list
GOVET:=$(GOCMD) vet
GOPATH:=$(shell $(GOCMD) env GOPATH)
BINARYVERSION:=$(shell git describe --tags)
ifndef BINARYVERSION
BINARYVERSION=$(shell git log -n 1 --pretty=format:"%h")
endif
GITLASTLOG=$(shell git log --pretty=format:'%h - %s (%cd) <%an>' -1)

GOBUILD=CGO_ENABLED=1 $(GOCMD) build -trimpath -ldflags '\
		-X "github.com/kayuii/go-fan-ctrl/version.BinaryVersion=${BINARYVERSION}" \
		-X "github.com/kayuii/go-fan-ctrl/version.GitLastLog=${GITLASTLOG}" \
		-w'

u := $(if $(update),-u)

BINARY_NAME:=fanctrl
GOFILES:=$(shell find . -name "*.go" -type f)

export GO111MODULE := on
export CGO_ENABLED := 1
export GOARCH := amd64
export GOOS := linux

all: test build

mini: test build-mini

.PHONY: build
build: deps
	$(GOBUILD) -o $(BINARY_NAME) ./cmd/fanctrl

.PHONY: build-mini
build-mini: deps
	$(GOBUILD) -ldflags "-s -w" -o $(BINARY_NAME)-mini ./cmd/fanctrl

.PHONY: build-static
build-static: deps
	CGO_ENABLED=0 $(GOBUILD) -ldflags '-linkmode "external" -extldflags "-static" -w -s ' -o $(BINARY_NAME)-static ./cmd/fanctrl

.PHONY: build-static2
build-static2: deps
	CGO_ENABLED=1 $(GOBUILD) -ldflags '-linkmode "external" -extldflags "-static" -w -s ' -o $(BINARY_NAME)-static ./cmd/chiacli

.PHONY: install
install: deps
	$(GOINSTALL) ./cmd/fanctrl

.PHONY: clean
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

.PHONY: deps
deps:
	go mod tidy

.PHONY: devel-deps
devel-deps:
	GO111MODULE=off $(GOGET) -v -u \
		golang.org/x/lint/golint

# .PHONY: lint
# lint: devel-deps
# 	for PKG in $(PACKAGES); do golint -set_exit_status $$PKG || exit 1; done;

# .PHONY: vet
# vet: deps devel-deps
# 	$(GOVET) $(PACKAGES)

.PHONY: fmt
fmt:
	$(GOFMT) -s -w $(GOFILES)

.PHONY: fmt-check
fmt-check:
	@diff=$$($(GOFMT) -s -d $(GOFILES)); \
	if [ -n "$$diff" ]; then \
		echo "Please run 'make fmt' and commit the result:"; \
		echo "$${diff}"; \
		exit 1; \
	fi;
