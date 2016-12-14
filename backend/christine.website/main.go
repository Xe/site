package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/gernest/front"
)

type Post struct {
	Title   string `json:"title"`
	Link    string `json:"link"`
	Summary string `json:"summary"`
	Date    string `json:"date"`
}

type Posts []*Post

func (p Posts) Len() int { return len(p) }
func (p Posts) Less(i, j int) bool {
	iDate, _ := time.Parse("2006-01-02", p[i].Date)
	jDate, _ := time.Parse("2006-01-02", p[j].Date)

	return iDate.Unix() < jDate.Unix()
}
func (p Posts) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

var posts Posts

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

	sort.Sort(sort.Reverse(posts))
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
