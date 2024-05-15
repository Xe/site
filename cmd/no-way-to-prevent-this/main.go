package main

import (
	"flag"
	"os"
	"text/template"
	"time"

	"github.com/go-faker/faker/v4"
)

var (
	date        = flag.String("date", time.Now().Format(time.DateOnly), "Date of the CVE")
	cPlusPlus   = flag.Bool("c++", false, "If true, the project is written in C++")
	cve         = flag.String("cve", "", "CVE number")
	cveLink     = flag.String("cve-link", "", "CVE link")
	project     = flag.String("project", "", "Project name")
	projectLink = flag.String("project-link", "", "Project link")
	summary     = flag.String("summary", "a memory safety vulnerability resulting in arbitrary code execution", "Summary of the CVE")
)

func main() {
	flag.Parse()

	os.MkdirAll("./lume/src/shitposts/no-way-to-prevent-this", 0755)
	fout, err := os.Create("./lume/src/shitposts/no-way-to-prevent-this/" + *cve + ".md")
	if err != nil {
		panic(err)
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

	tmpl := template.Must(template.New("article").Parse(articleTemplate))
	if err := tmpl.Execute(fout, data); err != nil {
		panic(err)
	}
}

const articleTemplate = `---
title: '"No way to prevent this" say users of only language where this regularly happens'
date: {{.Date}}
series: "no-way-to-prevent-this"
type: blog
hero:
  ai: "Photo by Andrea Piacquadio, source: Pexels"
  file: sad-business-man
  prompt: A forlorn business man resting his head on a brown wall next to a window.
---

In the hours following the release of [{{.CVE}}]({{.CVELink}}) for the project [{{.Project}}]({{.ProjectLink}}), site reliability workers
and systems administrators scrambled to desperately rebuild and patch all their systems to fix {{.Summary}}. This is due to the affected components being
written in C{{if .CPlusPlus}}++{{end}}, the only programming language where these vulnerabilities regularly happen. "This was a terrible tragedy, but sometimes
these things just happen and there's nothing anyone can do to stop them," said programmer {{.Name}}, echoing statements
expressed by hundreds of thousands of programmers who use the only language where 90% of the world's memory safety vulnerabilities have
occurred in the last 50 years, and whose projects are 20 times more likely to have security vulnerabilities. "It's a shame, but what can
we do? There really isn't anything we can do to prevent memory safety vulnerabilities from happening if the programmer doesn't want to
write their code in a robust manner." At press time, users of the only programming language in the world where these vulnerabilities
regularly happen once or twice per quarter for the last eight years were referring to themselves and their situation as "helpless."
`
