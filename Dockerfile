FROM xena/christine.website:1.1-28-g7733f4c
EXPOSE 5000
RUN apk add --no-cache bash
CMD /site/run.sh
