FROM xena/christine.website:1.1-5-g883104c
EXPOSE 5000
RUN apk add --no-cache bash
CMD /site/run.sh
