FROM xena/go:1.10 AS build
COPY . /root/go/src/github.com/Xe/site
RUN CGO_ENABLED=0 GOBIN=/root go install github.com/Xe/site/cmd/site

FROM xena/alpine
EXPOSE 5000
RUN apk add --no-cache bash
COPY --from=build /root/site /site/site
COPY ./static2 /site/static
COPY ./templates /site/templates
COPY ./blog /site/blog
COPY ./css /site/css
COPY ./run.sh /site/run.sh

HEALTHCHECK CMD curl --fail http://127.0.0.1:5000 || exit 1
CMD /site/run.sh
