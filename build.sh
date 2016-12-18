#!/bin/bash

set -x

export PATH="$PATH:/usr/local/go/bin:/usr/local/node/bin"
export CI="true"

(cd /site/frontend
 yes | npm install
 npm install -g bower
 yes 2 | bower install --allow-root
 npm run build
 rm -rf bower_components node_modules) &

(cd /site/backend/christine.website
 go build
 mv christine.website /usr/bin) &

wait
