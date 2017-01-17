#!/bin/bash

set -e
set -x

(cd frontend \
        && rm -rf node_modules bower_components \
        && npm install && npm run build \
        && asar pack static ../frontend.asar \
        && cd .. \
        && keybase sign -d -i ./frontend.asar -o ./frontend.asar.sig)

box box.rb
