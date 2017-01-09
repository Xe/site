#!/bin/bash

set -e
set -x

(cd frontend && rm -rf node_modules bower_components && npm install && npm run build && asar pack static ../frontend.asar)

box box.rb
