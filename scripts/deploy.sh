#!/bin/sh

set -e
set -x

docker build .
git push dokku master
