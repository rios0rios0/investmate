SCRIPTS_DIR ?= $(HOME)/Development/github.com/rios0rios0/pipelines
-include $(SCRIPTS_DIR)/makefiles/common.mk
-include $(SCRIPTS_DIR)/makefiles/golang.mk

.PHONY: build build-musl debug install run

build:
	rm -rf bin
	go build -o bin/investmate ./cmd/investmate
	strip -s bin/investmate

debug:
	rm -rf bin
	go build -gcflags "-N -l" -o bin/investmate ./cmd/investmate

build-musl:
	CGO_ENABLED=1 CC=musl-gcc go build \
		--ldflags 'linkmode external -extldflags="-static"' -o bin/investmate ./cmd/investmate
	strip -s bin/investmate

run:
	go run ./cmd/investmate

install:
	make build
	mkdir -p ~/.local/bin
	cp -v bin/investmate ~/.local/bin/investmate
