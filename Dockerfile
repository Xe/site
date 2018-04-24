FROM xena/christine.website:1.1-35-gd15c99e
EXPOSE 5000
RUN apk add --no-cache bash
CMD /site/run.sh
