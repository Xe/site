FROM xena/christine.website:1.1-12-g94167f7
EXPOSE 5000
RUN apk add --no-cache bash
CMD /site/run.sh
