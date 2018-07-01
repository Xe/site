FROM xena/go:1.10 AS build
COPY . /root/go/src/github.com/Xe/site
RUN GOBIN=/root go build github.com/Xe/site

FROM xena/alpine
EXPOSE 5000
RUN apk add --no-cache bash
COPY --from=build /root/site /site/site
COPY ./templates /site/templates
COPY ./blog /site/blog
COPY ./run.sh /site.sh
COPY ./static /site/static

HEALTHCHECK CMD curl --fail http://127.0.0.1:5000 || exit 1
CMD /site/run.sh
