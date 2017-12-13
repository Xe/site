FROM xena/christine.website:1.0-3-gd5cfa1b
EXPOSE 5000
RUN apk add --no-cache bash
CMD /site/run.sh
