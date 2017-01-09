#!/bin/sh

set -e
set -x

export PATH="$PATH:/usr/local/go/bin"
export CI="true"

(cd /site/backend/christine.website
 go get github.com/Xe/asarfs
 go get github.com/gernest/front
 go get layeh.com/asar
 go get gopkg.in/yaml.v2
 go get github.com/urfave/negroni

 go build -v
 mv christine.website /usr/bin) &

wait
