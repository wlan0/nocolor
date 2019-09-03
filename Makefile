PWD := $(shell pwd)
GOPATH := $(shell go env GOPATH)
LDFLAGS := $(shell echo "")

GOOS := $(shell go env GOOS)
GOOSALT ?= 'linux'
ifeq ($(GOOS),'darwin')
  GOOSALT = 'mac'
endif

BUILD_LDFLAGS := '$(LDFLAGS)'

all: build

build:
	@echo "building nocolor binary to ./nocolor"
	@GOPROXY=https://proxy.golang.org GO111MODULE=on GO_FLAGS="" CGO_ENABLED=0 go build -tags kqueue --ldflags $(BUILD_LDFLAGS)

build-linux:
	@echo "building nocolor-linux-amd64 to ./nocolor-linux-amd64"
	@GOOS=linux GOARCH=amd64 GOPROXY=https://proxy.golang.org GO111MODULE=on GO_FLAGS="" CGO_ENABLED=0 go build -o nocolor-linux-amd64 -tags kqueue --ldflags $(BUILD_LDFLAGS)

build-darwin:
	@echo "building nocolor-linux-amd64 to ./nocolor-darwin-amd64"
	@GOOS=darwin GOARCH=amd64 GOPROXY=https://proxy.golang.org GO111MODULE=on GO_FLAGS="" CGO_ENABLED=0 go build -o nocolor-darwin-amd64 -tags kqueue --ldflags $(BUILD_LDFLAGS)

build-windows:
	@echo "building nocolor-linux-amd64 to ./nocolor-windows-amd64"
	@GOOS=windows GOARCH=amd64 GOPROXY=https://proxy.golang.org GO111MODULE=on GO_FLAGS="" CGO_ENABLED=0 go build -o nocolor-windows-amd64 -tags kqueue --ldflags $(BUILD_LDFLAGS)

release: build-linux build-darwin build-windows
	@echo "releasing multi-platform nocolor binaries to releases/"
	@mkdir -p releases/
	@mv nocolor-linux-amd64	releases/nocolor-linux-amd64
	@mv nocolor-windows-amd64	releases/nocolor-windows-amd64
	@mv nocolor-darwin-amd64	releases/nocolor-darwin-amd64
	@tar cvzf nocolor-multiplatform.tar.gz releases/
