#!/bin/sh

set -e
set -x

export PATH="$PATH:/usr/local/go/bin"
export CI="true"

cd /site/backend/christine.website
go build -v
mv christine.website /usr/bin
