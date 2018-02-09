FROM xena/christine.website:1.1-20-gf1471b8
EXPOSE 5000
RUN apk add --no-cache bash
CMD /site/run.sh
