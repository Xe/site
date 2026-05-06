ARG ALPINE_VERSION=3.20
ARG GO_VERSION=1.26.2
ARG UBUNTU_VERSION=26.04

# Go toolchain bootstrapper
FROM golang:${GO_VERSION} AS go
ARG GO_VERSION=1.26.2

RUN CGO_ENABLED=0 go install golang.org/dl/go${GO_VERSION}@latest \
  && mkdir -p /app/bin \
  && mv /go/bin/go${GO_VERSION} /app/bin/go

# Iosevka for the resume
FROM --platform=${BUILDPLATFORM} alpine:${ALPINE_VERSION} AS fonts
ARG FONTS_VERSION=20250421
WORKDIR /fonts
RUN set -x \
  && apk add -U unzip ca-certificates curl \
  && curl -Lo iosevka.zip https://files.xeiaso.net/dl/iosevka-${FONTS_VERSION}.zip \
  && unzip -d /fonts iosevka.zip

# dhall-json for configuration building
FROM --platform=${BUILDPLATFORM} alpine:${ALPINE_VERSION} AS dhall-json
ARG DHALL_VERSION=1.42.2
ARG DHALL_JSON_VERSION=1.7.12
RUN mkdir -p /app
WORKDIR /app
RUN set -x \
  && apk add -U curl bzip2 ca-certificates \
  && curl -L -o dhall-linux.tar.bz2 https://github.com/dhall-lang/dhall-haskell/releases/download/${DHALL_VERSION}/dhall-json-${DHALL_JSON_VERSION}-x86_64-linux.tar.bz2 \
  && tar -xvjf dhall-linux.tar.bz2 \
  && mv bin/dhall-to-json .

# deno 
FROM alpine:${ALPINE_VERSION} AS deno
ARG DENO_VERSION=2.2.11
RUN mkdir -p /app
WORKDIR /app
RUN apk add -U curl unzip ca-certificates \
  && curl -sSLo deno.zip https://github.com/denoland/deno/releases/download/v${DENO_VERSION}/deno-$(uname -m)-unknown-linux-gnu.zip \
  && unzip deno.zip

# typst
FROM alpine:${ALPINE_VERSION} AS typst
ARG TYPST_VERSION=0.14.2
RUN mkdir -p /app
WORKDIR /app
RUN set -x \
  && apk add -U curl xz ca-certificates \
  && curl -sSLo typst.tar.xz https://github.com/typst/typst/releases/download/v${TYPST_VERSION}/typst-$(uname -m)-unknown-linux-musl.tar.xz \
  && tar xJf typst.tar.xz -C . \
  && mv typst-$(uname -m)-unknown-linux-musl/typst .

# runtime image
FROM ubuntu:${UBUNTU_VERSION} AS run

RUN mkdir -p /app
WORKDIR /app

RUN apt-get update \
  && apt-get -y install ca-certificates media-types nodejs npm

#RUN apk add -U ca-certificates deno typst mailcap
ENV TYPST_FONT_PATHS=/app/fonts
ENV GOPATH=/go

COPY --from=go /app/bin/go /usr/local/bin/go
COPY --from=fonts /fonts /app/fonts
COPY --from=dhall-json /app/dhall-to-json /usr/local/bin/dhall-to-json
COPY --from=deno /app/deno /usr/local/bin/deno
COPY --from=typst /app/typst /usr/local/bin/typst

RUN go download

LABEL org.opencontainers.image.source="https://tangled.org/xeiaso.net/site"
LABEL org.opencontainers.image.title="Xesite CI image"
LABEL org.opencontainers.image.description="Intermediate image with everything installed for Xesite CI"