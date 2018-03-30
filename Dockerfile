FROM xena/christine.website:1.1-22-gb845aa9
EXPOSE 5000
RUN apk add --no-cache bash
CMD /site/run.sh
