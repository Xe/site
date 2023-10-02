#!/usr/bin/env bash

set -ex

nix build .#docker
docker load < ./result
docker push ghcr.io/xe/site/bin
~/.fly/bin/fly deploy
