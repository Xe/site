VERSION 0.8
FROM alpine:edge
WORKDIR /app

deps:
    FROM golang:1.23-alpine
    WORKDIR /app

    COPY go.mod go.sum ./
    RUN go mod download

    SAVE ARTIFACT go.mod

fonts:
    FROM alpine:edge
    WORKDIR /fonts
    RUN apk add -U unzip ca-certificates curl \
     && curl -Lo iosevka.zip https://cdn.xeiaso.net/static/pkg/iosevka/ttf.zip \
     && unzip -d /fonts iosevka.zip

    SAVE ARTIFACT /fonts/ttf

dhall-json:
    FROM alpine:edge
    RUN apk add -U curl bzip2 ca-certificates \
     && curl -L -o dhall-linux.tar.bz2 https://github.com/dhall-lang/dhall-haskell/releases/download/1.42.0/dhall-json-1.7.12-x86_64-linux.tar.bz2 \
     && tar -xvjf dhall-linux.tar.bz2 \
     && mv bin/dhall-to-json .

    SAVE ARTIFACT dhall-to-json

build-patreon-saasproxy:
    FROM +deps
    COPY . .

    RUN --mount=type=cache,target=/root/.cache CGO_ENABLED=0 go build -gcflags "all=-N -l" -o patreon-saasproxy ./cmd/patreon-saasproxy

    SAVE ARTIFACT patreon-saasproxy

patreon-saasproxy:
    FROM alpine:edge
    WORKDIR /app

    COPY +build-patreon-saasproxy/patreon-saasproxy /app/patreon-saasproxy

    RUN apk add -U ca-certificates

    CMD ["./patreon-saasproxy"]

    LABEL org.opencontainers.image.source="https://github.com/Xe/site"

    SAVE IMAGE --push ghcr.io/xe/site/patreon:latest

build-xesite:
    FROM +deps
    COPY . .

    RUN --mount=type=cache,target=/root/.cache CGO_ENABLED=0 go build -gcflags "all=-N -l" -o xesite ./cmd/xesite

    SAVE ARTIFACT xesite

xesite:
    FROM alpine:edge
    WORKDIR /app

    RUN apk add -U ca-certificates deno typst mailcap
    ENV TYPST_FONT_PATHS=/app/fonts

    COPY +build-xesite/xesite /app/xesite
    COPY +fonts/ttf /app/fonts
    COPY +dhall-json/dhall-to-json /usr/local/bin/dhall-to-json

    CMD ["/app/xesite"]

    LABEL org.opencontainers.image.source="https://github.com/Xe/site"

    SAVE IMAGE --push ghcr.io/xe/site/bin:latest

all:
    BUILD --platform=linux/amd64 +xesite
    BUILD --platform=linux/amd64 +patreon-saasproxy 
