package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/Xe/asarfs"
	"github.com/Xe/ln"
	"github.com/gernest/front"
	"github.com/gorilla/feeds"
	"github.com/urfave/negroni"
)

// Post is a single post summary for the menu.
type Post struct {
	Title   string `json:"title"`
	Link    string `json:"link"`
	Summary string `json:"summary,omitifempty"`
	Body    string `json:"body, omitifempty"`
	Date    string `json:"date"`
}

// Posts implements sort.Interface for a slice of Post objects.
type Posts []*Post

func (p Posts) Len() int { return len(p) }
func (p Posts) Less(i, j int) bool {
	iDate, _ := time.Parse("2006-01-02", p[i].Date)
	jDate, _ := time.Parse("2006-01-02", p[j].Date)

	return iDate.Unix() < jDate.Unix()
}
func (p Posts) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

var (
	posts Posts
	rbody string
)

func init() {
	err := filepath.Walk("./blog/", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		fin, err := os.Open(path)
		if err != nil {
			return err
		}
		defer fin.Close()

		content, err := ioutil.ReadAll(fin)
		if err != nil {
			// handle error
		}

		m := front.NewMatter()
		m.Handle("---", front.YAMLHandler)
		front, _, err := m.Parse(bytes.NewReader(content))
		if err != nil {
			return err
		}

		sp := strings.Split(string(content), "\n")
		sp = sp[4:]
		data := strings.Join(sp, "\n")

		p := &Post{
			Title: front["title"].(string),
			Date:  front["date"].(string),
			Link:  strings.Split(path, ".")[0],
			Body:  data,
		}

		posts = append(posts, p)

		return nil
	})

	if err != nil {
		panic(err)
	}

	sort.Sort(sort.Reverse(posts))

	resume, err := ioutil.ReadFile("./static/resume/resume.md")
	if err != nil {
		panic(err)
	}

	rbody = string(resume)
}

func main() {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {})
	http.HandleFunc("/api/blog/posts", writeBlogPosts)
	http.HandleFunc("/api/blog/post", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		name := q.Get("name")

		if name == "" {
			goto fail
		}

		for _, p := range posts {
			if strings.HasSuffix(p.Link, name) {
				json.NewEncoder(w).Encode(p)
				return
			}
		}

	fail:
		http.Error(w, "Not Found", http.StatusNotFound)
	})
	http.HandleFunc("/api/resume", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(struct {
			Body string `json:"body"`
		}{
			Body: rbody,
		})
	})

	if os.Getenv("USE_ASAR") == "yes" {
		log.Println("serving site frontend from asar file")

		do404 := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Not found", http.StatusNotFound)
		})
		fe, err := asarfs.New("./frontend.asar", do404)
		if err != nil {
			log.Fatal("frontend: ", err)
		}

		http.Handle("/dist/", http.FileServer(fe))
	} else {
		log.Println("serving site frontend from filesystem")
		http.Handle("/dist/", http.FileServer(http.Dir("./frontend/static/")))
	}

	http.Handle("/static/", http.FileServer(http.Dir(".")))
	http.HandleFunc("/", writeIndexHTML)

	port := os.Getenv("PORT")
	if port == "" {
		port = "9090"
	}

	http.HandleFunc("/blog.rss", createFeed)
	http.HandleFunc("/blog.atom", createAtom)

	n := negroni.Classic()
	n.UseHandler(http.DefaultServeMux)

	log.Fatal(http.ListenAndServe(":"+port, n))
}

var bootTime = time.Now()

var feed = &feeds.Feed{
	Title:       "Christine Dodrill's Blog",
	Link:        &feeds.Link{Href: "https://christine.website/blog"},
	Description: "My blog posts and rants about various technology things.",
	Author:      &feeds.Author{Name: "Christine Dodrill", Email: "me@christine.website"},
	Created:     bootTime,
	Copyright:   "This work is copyright Christine Dodrill. My viewpoints are my own and not the view of any employer past, current or future.",
}

func init() {
	for _, item := range posts {
		itime, _ := time.Parse("2006-01-02", item.Date)
		feed.Items = append(feed.Items, &feeds.Item{
			Title:       item.Title,
			Link:        &feeds.Link{Href: "https://christine.website/" + item.Link},
			Description: item.Summary,
			Created:     itime,
		})
	}
}

// IncrediblySecureSalt *******
const IncrediblySecureSalt = "hunter2"

func createFeed(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/rss+xml")
	w.Header().Set("ETag", Hash(bootTime.String(), IncrediblySecureSalt))

	err := feed.WriteRss(w)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		ln.Error(err, ln.F{
			"remote_addr": r.RemoteAddr,
			"action":      "generating_rss",
			"uri":         r.RequestURI,
			"host":        r.Host,
		})
	}
}

func createAtom(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/atom+xml")
	w.Header().Set("ETag", Hash(bootTime.String(), IncrediblySecureSalt))

	err := feed.WriteAtom(w)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		ln.Error(err, ln.F{
			"remote_addr": r.RemoteAddr,
			"action":      "generating_rss",
			"uri":         r.RequestURI,
			"host":        r.Host,
		})
	}
}

func writeBlogPosts(w http.ResponseWriter, r *http.Request) {
	p := []interface{}{}
	for _, post := range posts {
		p = append(p, struct {
			Title   string `json:"title"`
			Link    string `json:"link"`
			Summary string `json:"summary,omitifempty"`
			Date    string `json:"date"`
		}{
			Title:   post.Title,
			Link:    post.Link,
			Summary: post.Summary,
			Date:    post.Date,
		})
	}
	json.NewEncoder(w).Encode(p)
}

func writeIndexHTML(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./frontend/static/dist/index.html")
}
