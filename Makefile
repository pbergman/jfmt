VERSION_REF := $(shell git rev-parse --short HEAD)
VERSION_NAME := $(shell git describe --all | sed "s/^heads\///")

build:
	go build -x \
		-o jfmt \
		-v \
		-ldflags '-w -s -extldflags "-static" -X main.version=$(VERSION_NAME) -X main.ref=$(VERSION_REF)' \
		main.go && \
	strip jfmt