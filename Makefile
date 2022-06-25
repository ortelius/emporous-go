GO := go

GO_BUILD_PACKAGES := ./cmd/...
GO_BUILD_BINDIR :=./bin
GIT_COMMIT := $(or $(SOURCE_GIT_COMMIT),$(shell git rev-parse --short HEAD))

GO_LD_EXTRAFLAGS :=-X github.com/uor-framework/client/cli.version="$(shell git tag | sort -V | tail -1)" \
				   -X github.com/uor-framework/client/cli.buildData="dev" \
				   -X github.com/uor-framework/client/cli.commit="$(GIT_COMMIT)" \
				   -X github.com/uor-framework/client/cli.buildDate="$(shell date -u +'%Y-%m-%dT%H:%M:%SZ')"

build:
	mkdir -p ${GO_BUILD_BINDIR}
	$(GO) build -o $(GO_BUILD_BINDIR)/client -ldflags="$(GO_LD_EXTRAFLAGS)" $(GO_BUILD_PACKAGES)
.PHONY: build

vendor:
	$(GO) mod tidy
	$(GO) mod verify
	$(GO) mod vendor
.PHONY: vendor

clean:
	@rm -rf ./$(GO_BUILD_BINDIR)/*
.PHONY: clean

test-unit:
	$(GO) test $(GO_BUILD_FLAGS) -coverprofile=coverage.out -race -count=1 ./...
.PHONY: test-unit

sanity: vendor format vet
	git diff --exit-code
.PHONY: sanity

format: 
	$(GO) fmt ./...
.PHONY: format

vet: 
	$(GO) vet ./...
.PHONY: vet

all: clean vendor test-unit build
.PHONY: all