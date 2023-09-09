package config

import (
	"fmt"
	"strings"
	"time"

	"go.jetpack.io/tyson"
	"xeiaso.net/internal/jsonfeed"
)

type Config struct {
	Authors            map[string]Author `json:"authors"`
	Characters         []Character       `json:"characters"`
	ClackSet           []string          `json:"clackSet"`
	ContactLinks       []Link            `json:"contactLinks"`
	DefaultAuthor      Author            `json:"defaultAuthor"`
	JobHistory         []Job             `json:"jobHistory"`
	NotableProjects    []Link            `json:"notableProjects"`
	Pronouns           []Pronoun         `json:"pronouns"`
	SeriesDescriptions map[string]string `json:"seriesDescriptions"`
	VODs               []VOD             `json:"vods"`
	WebMentionEndpoint string            `json:"webMentionEndpoint"`
	Redirects          map[string]string `json:"redirects"`
}

type Author struct {
	Name     string   `json:"name"`
	Handle   string   `json:"handle"`
	Image    *string  `json:"image,omitempty"`
	URL      *string  `json:"url,omitempty"`
	SameAs   []string `json:"sameAs,omitempty"`
	JobTitle *string  `json:"jobTitle,omitempty"`
	InSystem bool     `json:"inSystem"`
	Pronouns Pronoun  `json:"pronouns"`
}

func (a Author) ToJSONFeedAuthor() jsonfeed.Author {
	url := ""
	if a.URL != nil {
		url = *a.URL
	}
	avatar := ""
	if a.Image != nil {
		avatar = *a.Image
	}

	return jsonfeed.Author{
		Name:   a.Name,
		URL:    url,
		Avatar: avatar,
	}
}

type Character struct {
	Name        string   `json:"name"`
	StickerName string   `json:"stickerName"`
	DefaultPose string   `json:"defaultPose"`
	Description string   `json:"description"`
	Pronouns    Pronoun  `json:"pronouns"`
	Stickers    []string `json:"stickers"`
}

type Company struct {
	Name     string   `json:"name"`
	URL      *string  `json:"url,omitempty"`
	Tagline  string   `json:"tagline"`
	Location Location `json:"location"`
	Defunct  bool     `json:"defunct"`
}

type Date struct {
	time.Time
}

func (d *Date) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var dateStr string
	if err := unmarshal(&dateStr); err != nil {
		return err
	}
	if dateStr == "" {
		return nil
	}
	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return err
	}
	d.Time = t
	return nil
}

func (d *Date) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	if s == "null" {
		return nil
	}
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return err
	}
	d.Time = t
	return nil
}

func (d Date) MarshalJSON() ([]byte, error) {
	if d.Time.IsZero() {
		return []byte("null"), nil
	}
	return []byte(fmt.Sprintf("%q", d.Time.Format("2006-01-02"))), nil
}

type Job struct {
	Company        Company    `json:"company"`
	Title          string     `json:"title"`
	Contract       bool       `json:"contract"`
	StartDate      *Date      `json:"startDate"`
	EndDate        *Date      `json:"endDate,omitempty"`
	DaysWorked     int        `json:"daysWorked"`
	DaysBetween    int        `json:"daysBetween"`
	Salary         Salary     `json:"salary"`
	LeaveReason    *string    `json:"leaveReason,omitempty"`
	Locations      []Location `json:"locations"`
	Highlights     []string   `json:"highlights"`
	HideFromResume bool       `json:"hideFromResume"`
}

type Link struct {
	URL         string  `json:"url"`
	Title       string  `json:"title"`
	Description *string `json:"description,omitempty"`
}

type Location struct {
	City            string `json:"city"`
	StateOrProvince string `json:"stateOrProvince"`
	Country         string `json:"country"`
	Remote          bool   `json:"remote"`
}

type NagMessage struct {
	Name    string `json:"name"`
	Mood    string `json:"mood"`
	Message string `json:"message"`
}

type Person struct {
	Name  string   `json:"name"`
	Tags  []string `json:"tags"`
	Links []Link   `json:"links"`
}

type Pronoun struct {
	Nominative           string `json:"nominative"`
	Accusative           string `json:"accusative"`
	PossessiveDeterminer string `json:"possessiveDeterminer"`
	Possessive           string `json:"possessive"`
	Reflexive            string `json:"reflexive"`
	Singular             bool   `json:"singular"`
}

type Resume struct {
	Name                string   `json:"name"`
	Tagline             string   `json:"tagline"`
	Location            Location `json:"location"`
	Buzzwords           []string `json:"buzzwords"`
	Jobs                []Job    `json:"jobs"`
	NotablePublications []Link   `json:"notablePublications"`
}

type Salary struct {
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
	Per      string  `json:"per"`
	Stock    *Stock  `json:"stock,omitempty"`
}

type Stock struct {
	Kind         string `json:"kind"`
	Amount       int    `json:"amount"`
	Liquid       bool   `json:"liquid"`
	VestingYears int    `json:"vestingYears"`
	CliffYears   int    `json:"cliffYears"`
}

type VOD struct {
	Title       string   `json:"title"`
	Slug        string   `json:"slug"`
	Date        *Date    `json:"date"`
	Description string   `json:"description"`
	CDNPath     string   `json:"cdnPath"`
	Tags        []string `json:"tags"`
}

// Parse parses the config file and returns a Config struct.
func Parse(fname string) (*Config, error) {
	var c Config
	if err := tyson.Unmarshal(fname, &c); err != nil {
		return nil, err
	}
	return &c, nil
}
