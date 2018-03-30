FROM xena/christine.website:1.1-25-g6c29390
EXPOSE 5000
RUN apk add --no-cache bash
CMD /site/run.sh
