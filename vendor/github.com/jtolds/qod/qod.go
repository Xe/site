// Copyright (C) 2017 JT Olds
// See LICENSE for copying information.

// Package qod should NOT be used in a serious software engineering
// environment. qod stands for Quick and Dirty bahaha I just realized I got the
// acronym wrong. It's fine. It's on brand. Quick AND Dirty.
//
// The context is I noticed that Go is my favorite language, but when a task
// gets too complicated for a shell pipeline or awk or something, I turn to
// Python. Why not Go?
//
// In Python, I'd frequently write something like:
//
//   for line in sys.stdin:
//     vals = map(int, line.split())
//
// Here that is in Go:
//
//   package main
//
//   import (
//     "bufio"
//     "fmt"
//     "os"
//     "strconv"
//     "strings"
//   )
//
//   func main() {
//     scanner := bufio.NewScanner(os.Stdin)
//     for scanner.Scan() {
//       var vals []int64
//       for _, str := range strings.Fields(scanner.Text()) {
//         val, err := strconv.ParseInt(str, 10, 64)
//         if err != nil {
//           panic(err)
//         }
//         vals = append(vals, val)
//       }
//     }
//     if err := scanner.Err(); err != nil {
//       panic(err)
//     }
//   }
//
// Ugh! Considering I don't care about this throwaway shell pipeline
// replacement, I'm clearly fine with it blowing up if something's wrong, and
// wow this was too much.
//
// Package qod allows me to write the same type of thing in Go. Here is a
// reimplementation of the Python code above using qod:
//
//   package main
//
//   import (
//     "os"
//     "strings"
//
//     "github.com/jtolds/qod"
//   )
//
//   func main() {
//     for line := range qod.Lines(os.Stdin) {
//       vals := qod.Int64Slice(strings.Fields(line))
//     }
//   }
//
// Better! I'm more likely to use Go now for little scripts!
//
// Reminder: don't use this for anything real. Most of the stuff in here
// panics at the sight of any errors. That's obviously Bad and Wrong and you
// should actually handle your errors. Set up your build system's linter to
// reject anything that imports github.com/jtolds/qod please. If you have a
// build system for what you're doing at all this isn't for you. If you have
// some one-off tab-delimited data you need to process real quick like I seem
// to ALL THE TIME then okay.
package qod

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
	"unicode"
)

// ANE stands for Assert No Error. It panics if err != nil.
func ANE(err error) {
	if err != nil {
		panic(err)
	}
}

// AFH stands for Assert File Handle. It asserts there was no error and
// passes the file handle on through. Usage like:
//
//  fh := qod.AFH(os.Open(path))
func AFH(f *os.File, err error) *os.File {
	ANE(err)
	return f
}

// AI stands for Assert Int. It asserts there was no error and
// passes the int on through. Usage like:
//
//  qod.AI(fmt.Println("a line"))
func AI(i int, err error) int {
	ANE(err)
	return i
}

// Lines makes reading lines easier. Usage like:
//
//   for line := range Lines(os.Stdin) {
//     // do something with the line
//   }
//
// Returned lines will be right-stripped of whitespace.
// If you care about the lifetime of the channel that you're reading from and
// don't want it to leak, you probably shouldn't be using this package at all.
func Lines(r io.Reader) <-chan string {
	ch := make(chan string)
	go func() {
		defer close(ch)
		br := bufio.NewReader(r)
		for {
			l, err := br.ReadString('\n')
			if err == io.EOF {
				if l != "" {
					ch <- strings.TrimRightFunc(l, unicode.IsSpace)
				}
				break
			}
			ANE(err)
			ch <- strings.TrimRightFunc(l, unicode.IsSpace)
		}
	}()
	return ch
}

// Float64 converts a string to a float64
func Float64(val string) float64 {
	casted, err := strconv.ParseFloat(val, 64)
	ANE(err)
	return casted
}

// Float64Slice converts a []string to a []float64
func Float64Slice(vals []string) (rv []float64) {
	rv = make([]float64, 0, len(vals))
	for _, val := range vals {
		rv = append(rv, Float64(val))
	}
	return rv
}

// Int64 converts a string to an int64
func Int64(val string) int64 {
	casted, err := strconv.ParseInt(val, 10, 64)
	ANE(err)
	return casted
}

// Int64Slice converts a []string to an []int64
func Int64Slice(vals []string) (rv []int64) {
	rv = make([]int64, 0, len(vals))
	for _, val := range vals {
		rv = append(rv, Int64(val))
	}
	return rv
}

// Printlnf is just cause I constantly use Println, then turn it into Printf,
// then get frustrated I forgot the newline.
func Printlnf(format string, vals ...interface{}) {
	AI(fmt.Printf(format+"\n", vals...))
}

// Bytes will take an integer amount of bytes and format it with units.
func Bytes(amount int64) string {
	val := float64(amount)
	moves := 0
	for val >= 1024 {
		val /= 1024
		moves += 1
	}
	return fmt.Sprintf("%0.02f %s", val, []string{
		"B", "KiB", "MiB", "GiB", "TiB", "PiB", "EiB", "ZiB", "YiB"}[moves])
}

// SortedKeysBool returns the keys of a map[string]bool in sorted order.
func SortedKeysBool(v map[string]bool) []string {
	rv := make([]string, 0, len(v))
	for key := range v {
		rv = append(rv, key)
	}
	sort.Strings(rv)
	return rv
}
