FROM xena/christine.website:1.1-32-gd8c75d6
EXPOSE 5000
RUN apk add --no-cache bash
CMD /site/run.sh
