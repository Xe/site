FROM xena/christine.website:1.1-47-g3228e3b
EXPOSE 5000
RUN apk add --no-cache bash
CMD /site/run.sh
