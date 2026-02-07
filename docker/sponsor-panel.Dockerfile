ARG GO_VERSION=1.25
ARG ALPINE_VERSION=edge

FROM --platform=${BUILDPLATFORM} golang:${GO_VERSION}-alpine AS build

ARG TARGETOS
ARG TARGETARCH

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN apk -U add nodejs npm \
  && npm ci \
  && cd ./cmd/sponsor-panel \
  && go generate ./...

RUN --mount=type=cache,target=/root/.cache GOOS=${TARGETOS} GOARCH=${TARGETARCH} CGO_ENABLED=0 go build -gcflags "all=-N -l" -o /app/bin/sponsor-panel ./cmd/sponsor-panel

FROM alpine:${ALPINE_VERSION} AS run
WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY --from=build /app/bin/sponsor-panel /app/bin/sponsor-panel

EXPOSE 4823

CMD ["/app/bin/sponsor-panel"]

LABEL org.opencontainers.image.source="https://github.com/Xe/site"
LABEL org.opencontainers.image.title="Sponsor Panel Service"
LABEL org.opencontainers.image.description="Web panel for GitHub sponsors to manage their benefits"
