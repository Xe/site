FROM xena/christine.website:1.1-14-g74ad1ae
EXPOSE 5000
RUN apk add --no-cache bash
CMD /site/run.sh
