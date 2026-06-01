ARG GO_VERSION=1.26
ARG ALPINE_VERSION=edge

FROM --platform=${BUILDPLATFORM} golang:${GO_VERSION}-alpine AS build

ARG TARGETOS
ARG TARGETARCH

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN apk -U add git

RUN --mount=type=cache,target=/root/.cache GOOS=${TARGETOS} GOARCH=${TARGETARCH} CGO_ENABLED=0 go build -ldflags="-X xeiaso.net/v4.Version=$(git describe --tags --always --dirty)" -o /app/bin/futuresight ./cmd/futuresight

FROM alpine:${ALPINE_VERSION} AS run
WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY --from=build /app/bin/futuresight /app/bin/futuresight

EXPOSE 3000

CMD ["/app/bin/futuresight"]

LABEL org.opencontainers.image.source="https://tangled.org/xeiaso.net/site"
LABEL org.opencontainers.image.title="FutureSight Preview Service"
LABEL org.opencontainers.image.description="Serves content-addressed xesite preview builds from erofs volumes stored in Tigris"
