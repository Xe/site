# qod

[See the Documentation](https://godoc.org/github.com/jtolds/qod)

Package `qod` should NOT be used in a serious software engineering
environment. `qod` stands for Quick and Dirty bahaha I just realized I got the
acronym wrong. It's fine. It's on brand. Quick AND Dirty.

The context is I noticed that Go is my favorite language, but when a task
gets too complicated for a shell pipeline or `awk` or something, I turn to
Python. Why not Go?

In Python, I'd frequently write something like:

```python
for line in sys.stdin:
  vals = map(int, line.split())
```

Here that is in Go:

```go
package main

import (
  "bufio"
  "fmt"
  "os"
  "strconv"
  "strings"
)

func main() {
  scanner := bufio.NewScanner(os.Stdin)
  for scanner.Scan() {
    var vals []int64
    for _, str := range strings.Fields(scanner.Text()) {
      val, err := strconv.ParseInt(str, 10, 64)
      if err != nil {
        panic(err)
      }
      vals = append(vals, val)
    }
  }
  if err := scanner.Err(); err != nil {
    panic(err)
  }
}
```

Ugh! Considering I don't care about this throwaway shell pipeline
replacement, I'm clearly fine with it blowing up if something's wrong, and
wow this was too much.

`qod` allows me to write the same type of thing in Go. Here is a
reimplementation of the Python code above using `qod`:

```go
package main

import (
  "os"
  "strings"

  "github.com/jtolds/qod"
)

func main() {
  for line := range qod.Lines(os.Stdin) {
    vals := qod.Int64Slice(strings.Fields(line))
  }
}
```

Better! I'm more likely to use Go now for little scripts!

*Reminder:* don't use this for anything real. Most of the stuff in here
panics at the sight of any errors. That's obviously Bad and Wrong and you
should actually handle your errors. Set up your build system's linter to
reject anything that imports `github.com/jtolds/qod` please. If you have a
build system for what you're doing at all this isn't for you. If you have
some one-off tab-delimited data you need to process real quick like I seem
to ALL THE TIME then okay.

### License

Copyright (C) 2017 JT Olds. See LICENSE for copying information.
