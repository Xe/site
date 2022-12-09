#!/usr/bin/env bash

denobuild() {
    deno cache --import-map=./import_map.json --lock lock.json --lock-write *.tsx deps.ts
    deno bundle --import-map=./import_map.json --lock lock.json $1 $2
}

set -e

export RUST_LOG=info
denobuild ./mastodon_share_button.tsx ../../static/js/mastodon_share_button.js
