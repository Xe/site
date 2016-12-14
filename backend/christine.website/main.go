package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/gernest/front"
)

type Post struct {
	Title   string `json:"title"`
	Link    string `json:"link"`
	Summary string `json:"summary"`
	Date    string `json:"date"`
}

var posts []*Post

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

		m := front.NewMatter()
		m.Handle("---", front.YAMLHandler)
		front, _, err := m.Parse(fin)
		if err != nil {
			return err
		}

		p := &Post{
			Title: front["title"].(string),
			Date:  front["date"].(string),
			Link:  strings.Split(path, ".")[0],
		}

		posts = append(posts, p)

		return nil
	})

	if err != nil {
		panic(err)
	}
}

func main() {
	http.HandleFunc("/api/blog/posts", writeBlogPosts)
	http.HandleFunc("/api/blog/post", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		name := q.Get("name")

		fin, err := os.Open(path.Join("./blog", name))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer fin.Close()

		m := front.NewMatter()
		m.Handle("---", front.YAMLHandler)
		_, body, err := m.Parse(fin)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		fmt.Fprintln(w, body)
	})
	http.Handle("/dist/", http.FileServer(http.Dir("./frontend/static/")))
	http.HandleFunc("/", writeIndexHTML)

	log.Fatal(http.ListenAndServe(":9090", nil))
}

func writeBlogPosts(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(posts)
}

func writeIndexHTML(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./frontend/static/dist/index.html")
}
