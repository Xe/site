#!/usr/bin/env

inotifywait -m bin/ -m internal/embedded/generate.go -m blog/ -m talks/ -m gallery/ -e create -e moved_to |
    while read path action file; do
        echo "The file '$file' appeared in directory '$path' via '$action'"
        go generate -v ./...
    done
