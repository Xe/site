ARG ALPINE_VERSION=edge
ARG GO_VERSION=1.25
ARG UBUNTU_VERSION=24.04

FROM --platform=${BUILDPLATFORM} golang:${GO_VERSION} AS build

ARG TARGETOS
ARG TARGETARCH

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN --mount=type=cache,target=/root/.cache GOOS=${TARGETOS} GOARCH=${TARGETARCH} CGO_ENABLED=0 go build -gcflags "all=-N -l" -o /app/bin/xesite ./cmd/xesite

#
# Images for various facets of xesite
#

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
ARG TYPST_VERSION=0.13.1
RUN mkdir -p /app
WORKDIR /app
RUN set -x \
  && apk add -U curl xz ca-certificates \
  && curl -sSLo typst.tar.xz https://github.com/typst/typst/releases/download/v${TYPST_VERSION}/typst-$(uname -m)-unknown-linux-musl.tar.xz \
  && tar xJf typst.tar.xz -C . \
  && mv typst-$(uname -m)-unknown-linux-musl/typst .

# runtime image
FROM ubuntu:${UBUNTU_VERSION} AS run
WORKDIR /app

RUN apt-get update \
  && apt-get -y install ca-certificates media-types 

#RUN apk add -U ca-certificates deno typst mailcap
ENV TYPST_FONT_PATHS=/app/fonts

COPY --from=build /app/bin/xesite /app/bin/xesite
COPY --from=fonts /fonts /app/fonts
COPY --from=dhall-json /app/dhall-to-json /usr/local/bin/dhall-to-json
COPY --from=deno /app/deno /usr/local/bin/deno
COPY --from=typst /app/typst /usr/local/bin/typst

CMD ["/app/bin/xesite"]

LABEL org.opencontainers.image.source="https://github.com/Xe/site"
