#!/usr/bin/env bash

set -e

export RUST_LOG=info
deno cache --import-map=./import_map.json --lock deno.lock --lock-write **/*.tsx build.ts

deno run -A build.ts ./components/*.tsx
