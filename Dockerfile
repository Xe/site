FROM xena/go:1.12.1 AS build
ENV GOPROXY https://cache.greedo.xeserv.us
COPY . /site
WORKDIR /site
RUN CGO_ENABLED=0 GOBIN=/root go install -v ./cmd/site

FROM xena/alpine
EXPOSE 5000
RUN apk add --no-cache bash
WORKDIR /site
COPY --from=build /root/site .
COPY ./static /site/static
COPY ./templates /site/templates
COPY ./blog /site/blog
COPY ./css /site/css
COPY ./app /app
COPY ./app.json .
HEALTHCHECK CMD wget --spider http://127.0.0.1:5000 || exit 1
CMD ./site
