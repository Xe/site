#!/usr/bin/env bash

twirp-openapi-gen -in ./xeiaso/net/v1/xesite.proto -in ./within/website/x/mi/v1/mi.proto -out openapi.json -path-prefix /api -servers https://xeiaso.net -title "Xe Iaso Dot Net API" -doc-version 4.0.0
