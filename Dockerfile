FROM xena/christine.website:1.1-9-ga12904f
EXPOSE 5000
RUN apk add --no-cache bash
CMD /site/run.sh
