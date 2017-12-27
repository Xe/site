FROM xena/christine.website:1.1-3-g6197eca
EXPOSE 5000
RUN apk add --no-cache bash
CMD /site/run.sh
