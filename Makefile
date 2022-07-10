GO := go

GO_BUILD_PACKAGES := ./cmd/...
GO_BUILD_BINDIR :=./bin
GIT_COMMIT := $(or $(SOURCE_GIT_COMMIT),$(shell git rev-parse --short HEAD))
GIT_TAG :="$(shell git tag | sort -V | tail -1)"

GO_LD_EXTRAFLAGS :=-X github.com/uor-framework/client/cli.version="$(shell git tag | sort -V | tail -1)" \
				   -X github.com/uor-framework/client/cli.buildData="dev" \
				   -X github.com/uor-framework/client/cli.commit="$(GIT_COMMIT)" \
				   -X github.com/uor-framework/client/cli.buildDate="$(shell date -u +'%Y-%m-%dT%H:%M:%SZ')"

build:
	mkdir -p ${GO_BUILD_BINDIR}
	$(GO) build -o $(GO_BUILD_BINDIR)/$(ARCH)-client -ldflags="$(GO_LD_EXTRAFLAGS)" $(GO_BUILD_PACKAGES)
.PHONY: build

build-all-arch:
	@rm -rf ./$(GO_BUILD_BINDIR)/*
	env GOOS=linux   GOARCH=amd64	$(GO) build -o $(GO_BUILD_BINDIR)/client-linux-amd64   -ldflags="$(GO_LD_EXTRAFLAGS)" $(GO_BUILD_PACKAGES)
	env GOOS=linux   GOARCH=arm64   $(GO) build -o $(GO_BUILD_BINDIR)/client-linux-arm64   -ldflags="$(GO_LD_EXTRAFLAGS)" $(GO_BUILD_PACKAGES)
	env GOOS=linux   GOARCH=s390x	$(GO) build -o $(GO_BUILD_BINDIR)/client-linux-s390x   -ldflags="$(GO_LD_EXTRAFLAGS)" $(GO_BUILD_PACKAGES)
	env GOOS=linux   GOARCH=ppc64le $(GO) build -o $(GO_BUILD_BINDIR)/client-linux-ppc64le -ldflags="$(GO_LD_EXTRAFLAGS)" $(GO_BUILD_PACKAGES)
	env GOOS=linux   GOARCH=riscv64 $(GO) build -o $(GO_BUILD_BINDIR)/client-linux-riscv64 -ldflags="$(GO_LD_EXTRAFLAGS)" $(GO_BUILD_PACKAGES)
	env GOOS=darwin  GOARCH=amd64	$(GO) build -o $(GO_BUILD_BINDIR)/client-darwin-amd64  -ldflags="$(GO_LD_EXTRAFLAGS)" $(GO_BUILD_PACKAGES)
	env GOOS=darwin  GOARCH=arm64	$(GO) build -o $(GO_BUILD_BINDIR)/client-darwin-arm64  -ldflags="$(GO_LD_EXTRAFLAGS)" $(GO_BUILD_PACKAGES)
	env GOOS=windows GOARCH=amd64	$(GO) build -o $(GO_BUILD_BINDIR)/client-windows-amd64 -ldflags="$(GO_LD_EXTRAFLAGS)" $(GO_BUILD_PACKAGES)
.PHONY: multi-arch

container:
	buildah manifest create ghcr.io/uor-framework/client:$(GIT_TAG)
	buildah build --manifest ghcr.io/uor-framework/client --build-arg TARGETARCH=arm64   --arch arm64   --volume ${PWD}/bin:/data:z --tag ghcr.io/uor-framework/client-arm64:$(GIT_TAG)   --file ./Containerfile
	buildah build --manifest ghcr.io/uor-framework/client --build-arg TARGETARCH=amd64   --arch amd64   --volume ${PWD}/bin:/data:z --tag ghcr.io/uor-framework/client-amd64:$(GIT_TAG)   --file ./Containerfile
	buildah build --manifest ghcr.io/uor-framework/client --build-arg TARGETARCH=s390x   --arch s390x   --volume ${PWD}/bin:/data:z --tag ghcr.io/uor-framework/client-s390x:$(GIT_TAG)   --file ./Containerfile
	buildah build --manifest ghcr.io/uor-framework/client --build-arg TARGETARCH=ppc64le --arch ppc64le --volume ${PWD}/bin:/data:z --tag ghcr.io/uor-framework/client-ppc64le:$(GIT_TAG) --file ./Containerfile
.PHONY: container

container-push:
	buildah manifest push --all ghcr.io/uor-framework/client docker://ghcr.io/uor-framework/client:$(GIT_TAG)
.PHONY: container-push

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

all: clean vendor test-unit build-all-arch container
.PHONY: all
