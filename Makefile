entr := $(shell which entr)
upx := $(shell which upx)

run:
ifdef entr
	find . -name '*.go' 2>&1 | entr -r go run .
else
	go run .
endif

build:
	go build .
ifdef upx
	upx $$(basename $$(pwd))
endif

build-plugins:
	for p in $$(ls pkg/crawler/sources | grep -v .so) ; do\
	    go build \
			-buildmode=plugin \
			-o pkg/crawler/sources/$${p}.so \
			pkg/crawler/sources/$${p}/main.go;\
	done

migrate:
	go run ./cmd/db

test-dev:
	ENV=../../.env go test -v internal/database/*

.PHONY: run build migrate
