package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/facebookgo/flagenv"
	"github.com/go-git/go-git/v5"
	"xeiaso.net/v4/internal/lume"
)

var (
	bind         = flag.String("bind", ":3000", "Port to listen on")
	gitBranch    = flag.String("git-branch", "static_site", "Git branch to clone")
	gitRepo      = flag.String("git-repo", "https://github.com/Xe/site", "Git repository to clone")
	githubSecret = flag.String("github-secret", "", "GitHub secret to use for webhooks")
	siteURL      = flag.String("site-url", "https://kaine.shark-harmonic.ts.net/", "URL to use for the site")
)

func main() {
	flagenv.Parse()
	flag.Parse()

	ctx := context.Background()

	fs, err := lume.New(ctx, &lume.Options{
		Branch:        *gitBranch,
		Repo:          *gitRepo,
		StaticSiteDir: "lume",
		URL:           *siteURL,
	})
	if err != nil {
		log.Fatal(err)
	}

	defer fs.Close()

	if err != nil {
		log.Fatal(err)
	}

	http.Handle("/", http.FileServer(http.FS(fs)))

	http.HandleFunc("/.within/hook/github", func(w http.ResponseWriter, r *http.Request) {
		if err := fs.Update(r.Context()); err != nil {
			if err == git.NoErrAlreadyUpToDate {
				w.WriteHeader(http.StatusOK)
				fmt.Fprintln(w, "already up to date")
				return
			}
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	log.Printf("listening on %s", *bind)
	log.Fatal(http.ListenAndServe(*bind, nil))
}
