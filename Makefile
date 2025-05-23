PKG:=github.com/kayrus/gof5
APP_NAME:=gof5
PWD:=$(shell pwd)
UID:=$(shell id -u)
VERSION:=$(shell git describe --tags --always --dirty="-dev")
GOOS:=$(shell go env GOOS)
LDFLAGS:=-X main.Version=$(VERSION) -w -s
GOOS:=$(strip $(shell go env GOOS))
GOARCHs:=$(strip $(shell go env GOARCH))

ifeq "$(GOOS)" "windows"
SUFFIX=.exe
endif

# CGO must be enabled
export CGO_ENABLED:=1

build: fmt vet
	$(foreach GOARCH,$(GOARCHs),$(shell GOARCH=$(GOARCH) go build -ldflags="$(LDFLAGS)" -trimpath -o bin/$(APP_NAME)_$(GOOS)_$(GOARCH)$(SUFFIX) ./cmd/gof5))

docker:
	docker pull golang:latest
	docker run -ti --rm -e GOCACHE=/tmp -v $(PWD):/$(APP_NAME) -u $(UID):$(UID) --workdir /$(APP_NAME) golang:latest make

fmt:
	gofmt -s -w cmd pkg

vet:
	go vet ./...

static:
	staticcheck ./cmd/... ./pkg/...

test:
	go test -v ./cmd/... ./pkg/...
