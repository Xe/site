FROM xena/christine.website:1.1-39-gd3b28f5
EXPOSE 5000
RUN apk add --no-cache bash
CMD /site/run.sh
