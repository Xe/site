ARG GO_VERSION=1.25
ARG ALPINE_VERSION=edge

FROM --platform=${BUILDPLATFORM} golang:${GO_VERSION}-alpine AS build

ARG TARGETOS
ARG TARGETARCH

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN --mount=type=cache,target=/root/.cache GOOS=${TARGETOS} GOARCH=${TARGETARCH} CGO_ENABLED=0 go build -gcflags "all=-N -l" -o /app/bin/github-sponsor-webhook ./cmd/github-sponsor-webhook

FROM alpine:${ALPINE_VERSION} AS run
WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY --from=build /app/bin/github-sponsor-webhook /app/bin/github-sponsor-webhook

EXPOSE 8080

CMD ["/app/bin/github-sponsor-webhook"]

LABEL org.opencontainers.image.source="https://github.com/Xe/site"
LABEL org.opencontainers.image.title="GitHub Sponsors Webhook Service"
LABEL org.opencontainers.image.description="Standalone webhook service for processing GitHub Sponsors events"