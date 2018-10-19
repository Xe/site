FROM xena/go:1.11.1 AS build
ENV GOPROXY https://cache.greedo.xeserv.us
COPY . /site
WORKDIR /site
RUN CGO_ENABLED=0 GOBIN=/root go install -v ./cmd/site

FROM xena/alpine
EXPOSE 5000
RUN apk add --no-cache bash
COPY --from=build /root/site /site/site
COPY ./static /site/static
COPY ./templates /site/templates
COPY ./blog /site/blog
COPY ./css /site/css
COPY ./run.sh /site/run.sh
COPY ./app /app

HEALTHCHECK CMD wget --spider http://127.0.0.1:5000 || exit 1
CMD /site/run.sh
