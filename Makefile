GO := go

GO_BUILD_PACKAGES := ./cmd/...
GO_BUILD_BINDIR :=./bin
GIT_COMMIT := $(or $(SOURCE_GIT_COMMIT),$(shell git rev-parse --short HEAD))
GIT_TAG :="$(shell git tag | sort -V | tail -1)"

GO_LD_EXTRAFLAGS :=-X github.com/uor-framework/client/cli.version="$(shell git tag | sort -V | tail -1)" \
				   -X github.com/uor-framework/client/cli.buildData="dev" \
				   -X github.com/uor-framework/client/cli.commit="$(GIT_COMMIT)" \
				   -X github.com/uor-framework/client/cli.buildDate="$(shell date -u +'%Y-%m-%dT%H:%M:%SZ')"

build: prep-build-dir
	$(GO) build -o $(GO_BUILD_BINDIR)/client  -ldflags="$(GO_LD_EXTRAFLAGS)" $(GO_BUILD_PACKAGES)
.PHONY: build

cross-build-darwin-amd64:
	env GOOS=darwin  GOARCH=amd64	$(GO) build -o $(GO_BUILD_BINDIR)/client-darwin-amd64  -ldflags="$(GO_LD_EXTRAFLAGS)" $(GO_BUILD_PACKAGES)
.PHONY: cross-build-darwin-amd64

cross-build-darwin-arm64:
	env GOOS=darwin  GOARCH=arm64	$(GO) build -o $(GO_BUILD_BINDIR)/client-darwin-arm64  -ldflags="$(GO_LD_EXTRAFLAGS)" $(GO_BUILD_PACKAGES)
.PHONY: cross-build-darwin-arm64

cross-build-windows-amd64:
	env GOOS=windows GOARCH=amd64	$(GO) build -o $(GO_BUILD_BINDIR)/client-windows-amd64 -ldflags="$(GO_LD_EXTRAFLAGS)" $(GO_BUILD_PACKAGES)
.PHONY: cross-build-windows-amd64

cross-build-linux-amd64:
	env GOOS=linux   GOARCH=amd64	$(GO) build -o $(GO_BUILD_BINDIR)/client-linux-amd64   -ldflags="$(GO_LD_EXTRAFLAGS)" $(GO_BUILD_PACKAGES)
.PHONY: cross-build-linux-amd64

cross-build-linux-arm64:
	env GOOS=linux   GOARCH=arm64   $(GO) build -o $(GO_BUILD_BINDIR)/client-linux-arm64   -ldflags="$(GO_LD_EXTRAFLAGS)" $(GO_BUILD_PACKAGES)
.PHONY: cross-build-linux-arm64

cross-build-linux-ppc64le:
	env GOOS=linux   GOARCH=ppc64le $(GO) build -o $(GO_BUILD_BINDIR)/client-linux-ppc64le -ldflags="$(GO_LD_EXTRAFLAGS)" $(GO_BUILD_PACKAGES)
.PHONY: cross-build-linux-ppc64le

cross-build-linux-s390x:
	env GOOS=linux   GOARCH=s390x	$(GO) build -o $(GO_BUILD_BINDIR)/client-linux-s390x   -ldflags="$(GO_LD_EXTRAFLAGS)" $(GO_BUILD_PACKAGES)
.PHONY: cross-build-linux-s390x

cross-build-linux-riscv64:
	env GOOS=linux   GOARCH=riscv64 $(GO) build -o $(GO_BUILD_BINDIR)/client-linux-riscv64 -ldflags="$(GO_LD_EXTRAFLAGS)" $(GO_BUILD_PACKAGES)
.PHONY: cross-build-linux-riscv64

cross-build: prep-build-dir cross-build-darwin-amd64 cross-build-darwin-arm64 cross-build-windows-amd64 cross-build-linux-amd64 cross-build-linux-arm64 cross-build-linux-ppc64le cross-build-linux-s390x cross-build-linux-riscv64
.PHONY: cross-build

prep-build-dir:
	mkdir -p ${GO_BUILD_BINDIR}
.PHONY: prep-build-dir

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

cross-build-all: clean vendor test-unit cross-build
.PHONY: cross-build-all
