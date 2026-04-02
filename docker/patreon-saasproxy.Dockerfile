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
RUN --mount=type=cache,target=/root/.cache GOOS=${TARGETOS} GOARCH=${TARGETARCH} CGO_ENABLED=0 go build -ldflags="-X xeiaso.net/v4.Version=$(git describe --tags --always --dirty)"  -gcflags "all=-N -l" -o /app/bin/patreon-saasproxy ./cmd/patreon-saasproxy

FROM alpine:${ALPINE_VERSION} AS run
WORKDIR /app

COPY --from=build /app/bin/patreon-saasproxy /app/bin/patreon-saasproxy

CMD ["/app/bin/patreon-saasproxy"]

LABEL org.opencontainers.image.source="https://github.com/Xe/site"
