FROM xena/christine.website:1.1-37-gf103032
EXPOSE 5000
RUN apk add --no-cache bash
CMD /site/run.sh
