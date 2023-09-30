---
title: "Crazy Experiment: Ship the Frontend as an asar document"
date: "2017-01-09"
tags:
 - asar
 - frontend
---

Today's crazy experiment is using an [asar archive](https://github.com/electron/asar) for shipping around
and mounting frontend Javascript applications. This is something I feel is worth doing because it allows
the web frontend developer (or team) give the backend team a single "binary" that can be dropped into the
deployment process without having to build the frontend code as part of CI.

asar is an interesting file format because it allows for random access of the data inside the archive.
This allows an HTTP server to be able to concurrently serve files out of it without having to lock or
incur an additional open file descriptor.

In order to implement this, I have created a Go package named [asarfs](https://github.com/Xe/asarfs) that
exposes the contents of an asar archive as a standard http.Handler.

Example Usage:

```go
package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Xe/asarfs"
)

func do404(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not found", http.StatusNotFound)
}

func main() {
	fs, err := asarfs.New("./static.asar", http.HandlerFunc(do404))
	if err != nil {
		log.Fatal(err)
	}

	http.ListenAndServe(":"+os.Getenv("PORT"), fs)
}
```

I made some contrived benchmarks using some sample data (lots of large json files from mongodb dumps)
I had laying around and ran them a few times. The results were very promising:

```console
[~/g/s/g/X/asarfs] : go1.8beta2 test -bench=. -benchmem
BenchmarkHTTPFileSystem-8          20000             66481 ns/op            3219 B/op         58 allocs/op
BenchmarkASARfs-8                  20000             72084 ns/op            3549 B/op         77 allocs/op
BenchmarkPreloadedASARfs-8         20000             62894 ns/op            3218 B/op         58 allocs/op
PASS
ok      github.com/Xe/asarfs    5.636s
```

Amazingly, the performance and memory usage differences between serving the files over an asar archive
and off of the filesystem are negligible. I've implemented it in the latest release of my personal website
and hopefully end users should be seeing no difference in page load times.
