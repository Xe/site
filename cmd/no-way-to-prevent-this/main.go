package main

import (
	"embed"
	"flag"
	"log"
	"os"
	"text/template"
	"time"

	"github.com/go-faker/faker/v4"
)

var (
	date         = flag.String("date", time.Now().Format(time.DateOnly), "Date of the CVE")
	cPlusPlus    = flag.Bool("c++", false, "If true, the project is written in C++")
	cve          = flag.String("cve", "", "CVE number")
	cveLink      = flag.String("cve-link", "", "CVE link")
	project      = flag.String("project", "", "Project name")
	projectLink  = flag.String("project-link", "", "Project link")
	summary      = flag.String("summary", "a memory safety vulnerability resulting in arbitrary code execution", "Summary of the CVE")
	templateName = flag.String("template", "memory-safety", "Template name to use for the post")

	//go:embed templates/*.tmpl
	templates embed.FS
)

func main() {
	flag.Parse()

	tmpl, err := template.ParseFS(templates, "templates/*.tmpl")
	if err != nil {
		log.Fatalf("can't parse templates: %v", err)
	}

	os.MkdirAll("./lume/src/shitposts/no-way-to-prevent-this/"+*templateName, 0755)
	fout, err := os.Create("./lume/src/shitposts/no-way-to-prevent-this/" + *templateName + "/" + *cve + ".md")
	if err != nil {
		log.Fatalf("can't create output file: %v", err)
	}

	defer fout.Close()

	data := map[string]any{
		"Date":        *date,
		"CVE":         *cve,
		"CVELink":     *cveLink,
		"Project":     *project,
		"ProjectLink": *projectLink,
		"Summary":     *summary,
		"Name":        faker.Name(),
		"CPlusPlus":   *cPlusPlus,
	}

	t := tmpl.Lookup(*templateName + ".tmpl")
	if t == nil {
		log.Fatalf("can't find template %s", *templateName)
	}

	err = t.Execute(fout, data)
	if err != nil {
		log.Fatalf(
			"error writing template %s to %s: %v",
			*templateName,
			fout.Name(),
			err)
	}
}
