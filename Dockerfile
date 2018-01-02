FROM xena/christine.website:1.1-7-gbf4ccd5
EXPOSE 5000
RUN apk add --no-cache bash
CMD /site/run.sh
