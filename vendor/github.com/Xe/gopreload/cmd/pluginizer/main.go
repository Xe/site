package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "github.com/Xe/gopreload"
)

var (
	pkgName     = flag.String("pkg", "", "package to underscore import")
	destPkgName = flag.String("dest", "", "destination package to generate")
)

const codeTemplate = `//+build go1.8

package main

import _ "$PACKAGE_PATH"`

func main() {
	flag.Parse()

	if *pkgName == "" || *destPkgName == "" {
		log.Fatal("must set -pkg and -dest")
	}

	srcDir := filepath.Join(os.Getenv("GOPATH"), "src", *destPkgName)

	err := os.MkdirAll(srcDir, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	fout, err := os.Create(srcDir + "/main.go")
	if err != nil {
		log.Fatal(err)
	}
	defer fout.Close()

	codeBody := os.Expand(codeTemplate, func(s string) string {
		if s == "PACKAGE_PATH" {
			return *pkgName
		}

		return "no idea man"
	})

	fmt.Fprintln(fout, codeBody)

	fmt.Println("To build this plugin: ")
	fmt.Println("  $ go build -buildmode plugin -o /path/to/output.so " + *destPkgName)
}
