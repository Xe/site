FROM xena/go:1.11 AS build
COPY . /root/go/src/github.com/Xe/site
WORKDIR /root/go/src/github.com/Xe/site
RUN GO111MODULE=on CGO_ENABLED=0 GOBIN=/root go install -v -mod=vendor ./cmd/site

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
