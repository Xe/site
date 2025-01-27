package main

import (
	"embed"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
	"time"
)

var (
	//go:embed templates/*.tmpl
	templates embed.FS

	date = flag.String("date", time.Now().Format(time.DateOnly), "Publication date of the post")

	routing = map[string]string{
		"blog":     "lume/src/blog",
		"linkpost": "lume/src/blog",
		"note":     "lume/src/notes",
		"talk":     "lume/src/talks",
		"xecast":   "lume/src/xecast",
	}
)

// go run ./cmd/hydrate <kind> <slug>

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [flags] <kind> <slug>\n\n", filepath.Base(os.Args[0]))
		fmt.Fprintln(os.Stderr, "Available kinds:")

		templs, err := templates.ReadDir("templates")
		if err != nil {
			log.Panicf("can't read templates: %v", err)
		}

		for _, tmpl := range templs {
			kind, ok := strings.CutSuffix(tmpl.Name(), filepath.Ext(tmpl.Name()))
			if !ok {
				log.Panicf("can't cut extension from %q", tmpl.Name())
			}
			fmt.Fprintln(os.Stderr, "  *", kind)
		}

		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "Flags:")

		flag.PrintDefaults()

		os.Exit(2)
	}
}

func main() {
	tmpl, err := template.ParseFS(templates, "templates/*.tmpl")
	if err != nil {
		log.Fatalf("can't parse templates: %v", err)
	}

	flag.Parse()

	if flag.NArg() != 2 {
		flag.Usage()
	}

	kind := flag.Arg(0)
	slug := flag.Arg(1)

	year, err := yearOf(*date)
	if err != nil {
		log.Fatalf("can't parse year in %s: %v", *date, year)
	}

	os.MkdirAll(filepath.Join(routing[kind], year), 0755)

	foutName := filepath.Join(routing[kind], year, slug+".mdx")

	if _, err := os.Stat(foutName); !os.IsNotExist(err) {
		log.Printf("Potential error when trying to verify %s doesn't exist: %v", foutName, err)
		log.Println("Does the file already exist?")
		os.Exit(1)
	}

	fout, err := os.Create(foutName)
	if err != nil {
		log.Fatalf("can't create %s: %v", foutName, err)
	}

	if err := tmpl.ExecuteTemplate(fout, kind+".tmpl", struct {
		Date string
		Year string
		Slug string
	}{
		Date: *date,
		Year: year,
		Slug: slug,
	}); err != nil {
		log.Fatalf("error writing template: %v", err)
	}

	if err := fout.Close(); err != nil {
		log.Fatalf("error closing output file: %v", err)
	}

	codePath, err := exec.LookPath("code")
	if err != nil {
		log.Println("hint: control shift p -> install code command")
		log.Fatalf("can't find code command in $PATH: %v", err)
	}

	if err := exec.Command(codePath, foutName).Run(); err != nil {
		log.Fatalf("can't open %s in VS Code: %v", foutName, err)
	}

}

func yearOf(date string) (string, error) {
	t, err := time.Parse(time.DateOnly, date)
	if err != nil {
		return "", err
	}

	return strconv.Itoa(t.Year()), nil
}
