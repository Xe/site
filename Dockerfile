FROM xena/christine.website:1.1-41-g626ed4f
EXPOSE 5000
RUN apk add --no-cache bash
CMD /site/run.sh
