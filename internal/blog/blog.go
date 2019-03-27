package blog

import (
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"christine.website/internal/front"
	"github.com/russross/blackfriday"
)

// Post is a single blogpost.
type Post struct {
	Title    string        `json:"title"`
	Link     string        `json:"link"`
	Summary  string        `json:"summary,omitifempty"`
	Body     string        `json:"-"`
	BodyHTML template.HTML `json:"body"`
	Date     time.Time     `json:"date"`
}

// Posts implements sort.Interface for a slice of Post objects.
type Posts []Post

func (p Posts) Len() int { return len(p) }
func (p Posts) Less(i, j int) bool {
	iDate := p[i].Date
	jDate := p[j].Date

	return iDate.Unix() < jDate.Unix()
}
func (p Posts) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

// LoadPosts loads posts for a given directory.
func LoadPosts(path string) (Posts, error) {
	type postFM struct {
		Title string
		Date  string
	}
	var result Posts

	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
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
			return err
		}

		var fm postFM
		remaining, err := front.Unmarshal(content, &fm)
		if err != nil {
			return err
		}

		output := blackfriday.Run(remaining)

		const timeFormat = `2006-01-02`
		date, err := time.Parse(timeFormat, fm.Date)
		if err != nil {
			return err
		}

		p := Post{
			Title:    fm.Title,
			Date:     date,
			Link:     strings.Split(path, ".")[0],
			Body:     string(remaining),
			BodyHTML: template.HTML(output),
		}
		result = append(result, p)

		return nil
	})
	if err != nil {
		return nil, err
	}

	sort.Sort(sort.Reverse(result))

	return result, nil
}
