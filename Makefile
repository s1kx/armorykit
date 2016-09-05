APPNAME=armory-kit
SOURCEDIRS := $(glide novendor)
SOURCES := $(shell find $(SOURCEDIRS) -name '*.go')

VERSION=$(shell git describe)
BUILD_TIME=$(shell date +%FT%T%z)

BUILD_FLAGS=-X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME}

# Debug flags
GCFLAGS_DEBUG=-gcflags ""
LDFLAGS_DEBUG=-ldflags "${BUILD_FLAGS}"
# Release flags
GCFLAGS_RELEASE=-gcflags "-trimpath ${GOPATH}"
LDFLAGS_RELEASE=-ldflags "-s -w ${BUILD_FLAGS}"

.DEFAULT_GOAL: all

all: debug release

debug: $(SOURCES)
	go build ${GCFLAGS_DEBUG} ${LDFLAGS_DEBUG} -o bin/debug/${APPNAME}

release: $(SOURCES)
	go build ${GCFLAGS_RELEASE} ${LDFLAGS_RELEASE} -a -o bin/release/${APPNAME}
	strip bin/release/${APPNAME}

.PHONY: install
install:
	go install ${LDFLAGS}

.PHONY: test
test:
	go test $(shell glide novendor)
