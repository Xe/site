ARG GO_VERSION=1.24
ARG ALPINE_VERSION=edge
FROM --platform=${BUILDPLATFORM} golang:${GO_VERSION}-alpine AS build

ARG TARGETOS
ARG TARGETARCH

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN --mount=type=cache,target=/root/.cache GOOS=${TARGETOS} GOARCH=${TARGETARCH} CGO_ENABLED=0 go build -gcflags "all=-N -l" -o /app/bin/xesite ./cmd/xesite

FROM --platform=${BUILDPLATFORM} alpine:${ALPINE_VERSION} AS fonts

WORKDIR /fonts
RUN apk add -U unzip ca-certificates curl \
  && curl -Lo iosevka.zip https://cdn.xeiaso.net/static/pkg/iosevka/ttf.zip \
  && unzip -d /fonts iosevka.zip

FROM --platform=${BUILDPLATFORM} alpine:${ALPINE_VERSION} AS dhall-json

RUN mkdir -p /app
WORKDIR /app
RUN apk add -U curl bzip2 ca-certificates \
  && curl -L -o dhall-linux.tar.bz2 https://github.com/dhall-lang/dhall-haskell/releases/download/1.42.0/dhall-json-1.7.12-x86_64-linux.tar.bz2 \
  && tar -xvjf dhall-linux.tar.bz2 \
  && mv bin/dhall-to-json .

FROM alpine:${ALPINE_VERSION} AS run
WORKDIR /app

RUN apk add -U ca-certificates deno typst mailcap
ENV TYPST_FONT_PATHS=/app/fonts

COPY --from=build /app/bin/xesite /app/bin/xesite
COPY --from=fonts /fonts /app/fonts
COPY --from=dhall-json /app/dhall-to-json /usr/local/bin/dhall-to-json

CMD ["/app/bin/xesite"]

LABEL org.opencontainers.image.source="https://github.com/Xe/site"
