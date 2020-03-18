FROM xena/go:1.14 AS build
ENV GOPROXY https://cache.greedo.xeserv.us
COPY . /site
WORKDIR /site
RUN CGO_ENABLED=0 go test -v ./...
RUN CGO_ENABLED=0 GOBIN=/root go install -v ./cmd/site

FROM xena/alpine
EXPOSE 5000
WORKDIR /site
COPY --from=build /root/site .
COPY ./static /site/static
COPY ./templates /site/templates
COPY ./blog /site/blog
COPY ./talks /site/talks
COPY ./gallery /site/gallery
COPY ./css /site/css
COPY ./signalboost.dhall /site/signalboost.dhall
HEALTHCHECK CMD wget --spider http://127.0.0.1:5000/.within/health || exit 1
CMD ./site
