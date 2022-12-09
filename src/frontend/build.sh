#!/usr/bin/env bash

denobuild() {
    deno cache --import-map=./import_map.json --lock deno.lock --lock-write *.tsx deps.ts
    deno bundle --import-map=./import_map.json --lock deno.lock $1 $2
}

set -e

export RUST_LOG=info
denobuild ./mastodon_share_button.tsx ../../static/js/mastodon_share_button.js
denobuild ./wasiterm.tsx ../../static/js/wasiterm.js
