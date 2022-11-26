#!/usr/bin/env bash

set -e

export RUST_LOG=info
deno bundle ./mastodon_share_button.tsx ../../static/js/mastodon_share_button.js
