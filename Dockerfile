FROM xena/christine.website:1.1-16-g02a0bc0
EXPOSE 5000
RUN apk add --no-cache bash
CMD /site/run.sh
