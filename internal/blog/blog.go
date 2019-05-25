package blog

import (
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"christine.website/v1/internal/front"
	"github.com/russross/blackfriday"
)

// Post is a single blogpost.
type Post struct {
	Title      string        `json:"title"`
	Link       string        `json:"link"`
	Summary    string        `json:"summary,omitifempty"`
	Body       string        `json:"-"`
	BodyHTML   template.HTML `json:"body"`
	SlidesLink string        `json:"slides_link"`
	Date       time.Time
	DateString string `json:"date"`
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
func LoadPosts(path string, prepend string) (Posts, error) {
	type postFM struct {
		Title      string
		Date       string
		SlidesLink string `yaml:"slides_link"`
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

		fname := filepath.Base(path)
		fname = strings.TrimSuffix(fname, filepath.Ext(fname))

		p := Post{
			Title:      fm.Title,
			Date:       date,
			DateString: fm.Date,
			Link:       filepath.Join(prepend, fname),
			Body:       string(remaining),
			BodyHTML:   template.HTML(output),
			SlidesLink: fm.SlidesLink,
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
