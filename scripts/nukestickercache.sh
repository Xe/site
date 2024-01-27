#!/usr/bin/env bash

XEDNS=$(tailscale status --json | jq '.Peer | to_entries[] | .value.HostName | select(. | test("^xedn-[a-z]{3}$"))' -c -r | sort)
IFS=$'\n'

jo -a $*

for xedn in ${XEDNS}; do
    curl "http://${xedn}/xedn/purge" --data-binary "$(jo -a $*)" &
done

wait
