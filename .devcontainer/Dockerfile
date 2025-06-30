FROM alpine:3.21 AS dhall-json

ARG DHALL_VERSION=1.42.2
ARG DHALL_JSON_VERSION=1.7.12
RUN mkdir -p /app
WORKDIR /app
RUN set -x \
  && apk add -U curl bzip2 ca-certificates \
  && curl -L -o dhall-linux.tar.bz2 https://github.com/dhall-lang/dhall-haskell/releases/download/${DHALL_VERSION}/dhall-json-${DHALL_JSON_VERSION}-x86_64-linux.tar.bz2 \
  && tar -xvjf dhall-linux.tar.bz2 \
  && mv bin/dhall-to-json .

FROM ghcr.io/xe/devcontainer-base/pre/go

COPY --from=dhall-json /app/dhall-to-json /usr/local/bin/dhall-to-json