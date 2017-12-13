FROM xena/christine.website:dev
EXPOSE 5000
RUN apk add --no-cache bash
CMD /site/run.sh
