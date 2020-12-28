entr := $(shell which entr)
upx := $(shell which upx)

run:
ifdef entr
	find . -name '*.go' 2>&1 | entr -r go run .
else
	go run .
endif

build: build-api build-plugins

build-api:
	go build -ldflags "-s -w" .
ifdef upx
	upx $$(basename $$(pwd))
endif

build-plugins:
	for p in $$(ls pkg/crawler/sources | grep -v .so) ; do\
	    go build \
	    	-ldflags "-s -w" \
			-buildmode=plugin \
			-o pkg/crawler/sources/$${p}.so \
			pkg/crawler/sources/$${p}/main.go;\
	done

# TODO: Rename to db and add more helper functions
migrate:
	go run ./cmd/db
migrate-build:
	go build -ldflags "-s -w" -o db-helper ./cmd/db
ifdef upx
	upx db-helper
endif

test-dev:
	ENV=../../.env go test -v $$m

.PHONY: run build migrate
