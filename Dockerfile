FROM xena/christine.website:1.1-45-ge9d305a
EXPOSE 5000
RUN apk add --no-cache bash
CMD /site/run.sh
