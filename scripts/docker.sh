#!/bin/sh

set -e

docker build -t xena/site .
exec docker run --rm -itp 5030:5000 -e PORT=5000 xena/site
