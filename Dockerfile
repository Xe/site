FROM xena/christine.website:1.1-18-g4fb5948
EXPOSE 5000
RUN apk add --no-cache bash
CMD /site/run.sh
