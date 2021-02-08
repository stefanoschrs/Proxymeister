############################
# STEP 1 build executable binary
############################
FROM golang AS builder

RUN apt update -y \
    && apt install -y git make \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /build

COPY . .

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" .
RUN for p in $(ls pkg/crawler/sources | grep -v .so) ; do\
        CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build \
            -ldflags="-w -s" \
            -buildmode=plugin \
            -o pkg/crawler/sources/${p}.so \
            pkg/crawler/sources/${p}/main.go;\
    done

############################
# STEP 2 build a small image
############################
FROM golang

COPY --from=builder /build/proxymeister /proxymeister
COPY --from=builder /build/pkg/crawler/sources/*.so /

ENTRYPOINT ["/proxymeister"]
