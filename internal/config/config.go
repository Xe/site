package config

import (
	"encoding/json"
	"os/exec"
)

var (
	dhallToJSONBin string
)

func init() {
	var err error
	dhallToJSONBin, err = exec.LookPath("dhall-to-json")
	if err != nil {
		panic(err)
	}
}

func Load(fname string) (*Config, error) {
	cmd := exec.Command(dhallToJSONBin, "--file", fname)

	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var result Config
	if err := json.Unmarshal(out, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

type Config struct {
	Authors            map[string]Author    `json:"authors"`
	Characters         []Character          `json:"characters"`
	ClackSet           []string             `json:"clackSet"`
	ContactLinks       []Link               `json:"contactLinks"`
	DefaultAuthor      Author               `json:"defaultAuthor"`
	JobHistory         []Job                `json:"jobHistory"`
	MiToken            string               `json:"miToken"`
	NotableProjects    []NotableProjects    `json:"notableProjects"`
	Port               int                  `json:"port"`
	Pronouns           []Pronouns           `json:"pronouns"`
	SeriesDescMap      map[string]string    `json:"seriesDescMap"`
	SeriesDescriptions []SeriesDescriptions `json:"seriesDescriptions"`
	Signalboost        []Signalboost        `json:"signalboost"`
	WebMentionEndpoint string               `json:"webMentionEndpoint"`
	Resume             Resume               `json:"resume"`
}

type Pronouns struct {
	Accusative           string `json:"accusative"`
	Nominative           string `json:"nominative"`
	Possessive           string `json:"possessive"`
	PossessiveDeterminer string `json:"possessiveDeterminer"`
	Reflexive            string `json:"reflexive"`
	Singular             bool   `json:"singular"`
}

type Character struct {
	DefaultPose string   `json:"defaultPose"`
	Description string   `json:"description"`
	Name        string   `json:"name"`
	Pronouns    Pronouns `json:"pronouns"`
	StickerName string   `json:"stickerName"`
	Stickers    []string `json:"stickers"`
}

type Link struct {
	Description string `json:"description"`
	Title       string `json:"title"`
	URL         string `json:"url"`
}

type Author struct {
	Handle   string   `json:"handle"`
	Image    string   `json:"image"`
	InSystem bool     `json:"inSystem"`
	JobTitle string   `json:"jobTitle"`
	Name     string   `json:"name"`
	Pronouns Pronouns `json:"pronouns"`
	SameAs   []string `json:"sameAs"`
	URL      string   `json:"url"`
}

type Location struct {
	City            string `json:"city"`
	Country         string `json:"country"`
	Remote          bool   `json:"remote"`
	StateOrProvince string `json:"stateOrProvince"`
}

type Company struct {
	Defunct  bool     `json:"defunct"`
	Location Location `json:"location"`
	Name     string   `json:"name"`
	Tagline  string   `json:"tagline"`
	URL      string   `json:"url"`
}

type Salary struct {
	Amount   int    `json:"amount"`
	Currency string `json:"currency"`
	Per      string `json:"per"`
	Stock    *Stock `json:"stoc,omitempty"`
}

type Stock struct {
	Amount       int    `json:"amount"`
	CliffYears   int    `json:"cliffYears"`
	Kind         string `json:"kind"`
	Liquid       bool   `json:"liquid"`
	VestingYears int    `json:"vestingYears"`
}

type Job struct {
	Company        Company    `json:"company,omitempty"`
	Contract       bool       `json:"contract"`
	HideFromResume bool       `json:"hideFromResume"`
	Highlights     []string   `json:"highlights"`
	Locations      []Location `json:"locations"`
	Salary         Salary     `json:"salary,omitempty"`
	StartDate      string     `json:"startDate"`
	Title          string     `json:"title"`
	DaysWorked     int        `json:"daysWorked,omitempty"`
	EndDate        string     `json:"endDate,omitempty"`
	LeaveReason    string     `json:"leaveReason,omitempty"`
	DaysBetween    int        `json:"daysBetween,omitempty"`
}

type NotableProjects struct {
	Description string `json:"description"`
	Title       string `json:"title"`
	URL         string `json:"url"`
}

type SeriesDescriptions struct {
	Details string `json:"details"`
	Name    string `json:"name"`
}

type Signalboost struct {
	Links []Link   `json:"links"`
	Name  string   `json:"name"`
	Tags  []string `json:"tags"`
}

type Resume struct {
	Buzzwords           []string `json:"buzzwords"`
	Location            Location `json:"location"`
	Name                string   `json:"name"`
	NotablePublications []Link   `json:"notablePublications"`
	Tagline             string   `json:"tagline"`
}
