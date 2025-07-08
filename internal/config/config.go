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
	Characters         []Character          `json:"characters" dhall:"characters"`
	ClackSet           []string             `json:"clackSet" dhall:"clackSet"`
	ContactLinks       []Link               `json:"contactLinks" dhall:"contactLinks"`
	JobHistory         []Job                `json:"jobHistory" dhall:"jobHistory"`
	NotableProjects    []NotableProjects    `json:"notableProjects" dhall:"notableProjects"`
	Port               int                  `json:"port" dhall:"port"`
	Pronouns           []Pronouns           `json:"pronouns" dhall:"pronouns"`
	SeriesDescriptions []SeriesDescriptions `json:"seriesDescriptions" dhall:"seriesDescriptions"`
	Signalboost        []Signalboost        `json:"signalboost" dhall:"signalboost"`
	WebMentionEndpoint string               `json:"webMentionEndpoint" dhall:"webMentionEndpoint"`
	Resume             Resume               `json:"resume" dhall:"resume"`
}

type Pronouns struct {
	Accusative           string `json:"accusative" dhall:"accusative"`
	Nominative           string `json:"nominative" dhall:"nominative"`
	Possessive           string `json:"possessive" dhall:"possessive"`
	PossessiveDeterminer string `json:"possessiveDeterminer" dhall:"possessiveDeterminer"`
	Reflexive            string `json:"reflexive" dhall:"reflexive"`
	Singular             bool   `json:"singular" dhall:"singular"`
}

type Character struct {
	DefaultPose string   `json:"defaultPose" dhall:"defaultPose"`
	Description string   `json:"description" dhall:"description"`
	Name        string   `json:"name" dhall:"name"`
	Pronouns    Pronouns `json:"pronouns" dhall:"pronouns"`
	StickerName string   `json:"stickerName" dhall:"stickerName"`
	Stickers    []string `json:"stickers" dhall:"stickers"`
}

type Link struct {
	Description string `json:"description" dhall:"description"`
	Title       string `json:"title" dhall:"title"`
	URL         string `json:"url" dhall:"url"`
}

type Location struct {
	City            string `json:"city" dhall:"city"`
	Country         string `json:"country" dhall:"country"`
	Remote          bool   `json:"remote" dhall:"remote"`
	StateOrProvince string `json:"stateOrProvince" dhall:"stateOrProvince"`
}

type Company struct {
	Defunct  bool     `json:"defunct" dhall:"defunct"`
	Location Location `json:"location" dhall:"location"`
	Name     string   `json:"name" dhall:"name"`
	Tagline  string   `json:"tagline" dhall:"tagline"`
	URL      string   `json:"url" dhall:"url"`
}

type Salary struct {
	Amount   int    `json:"amount" dhall:"amount"`
	Currency string `json:"currency" dhall:"currency"`
	Per      string `json:"per" dhall:"per"`
	Stock    *Stock `json:"stoc,omitempty" dhall:"stoc,omitempty"`
}

type Stock struct {
	Amount       int    `json:"amount" dhall:"amount"`
	CliffYears   int    `json:"cliffYears" dhall:"cliffYears"`
	Kind         string `json:"kind" dhall:"kind"`
	Liquid       bool   `json:"liquid" dhall:"liquid"`
	VestingYears int    `json:"vestingYears" dhall:"vestingYears"`
}

type Job struct {
	Company        Company    `json:"company,omitempty" dhall:"company,omitempty"`
	Contract       bool       `json:"contract" dhall:"contract"`
	HideFromResume bool       `json:"hideFromResume" dhall:"hideFromResume"`
	Highlights     []string   `json:"highlights" dhall:"highlights"`
	Locations      []Location `json:"locations" dhall:"locations"`
	Salary         Salary     `json:"salary,omitempty" dhall:"salary,omitempty"`
	StartDate      string     `json:"startDate" dhall:"startDate"`
	Title          string     `json:"title" dhall:"title"`
	DaysWorked     int        `json:"daysWorked,omitempty" dhall:"daysWorked,omitempty"`
	EndDate        string     `json:"endDate,omitempty" dhall:"endDate,omitempty"`
	LeaveReason    string     `json:"leaveReason,omitempty" dhall:"leaveReason,omitempty"`
	DaysBetween    int        `json:"daysBetween,omitempty" dhall:"daysBetween,omitempty"`
}

type NotableProjects struct {
	Description string `json:"description" dhall:"description"`
	Title       string `json:"title" dhall:"title"`
	URL         string `json:"url" dhall:"url"`
}

type SeriesDescriptions struct {
	Details string `json:"details" dhall:"details"`
	Name    string `json:"name" dhall:"name"`
}

type Signalboost struct {
	Links []Link   `json:"links" dhall:"links"`
	Name  string   `json:"name" dhall:"name"`
	Tags  []string `json:"tags" dhall:"tags"`
}

type Resume struct {
	Buzzwords           []string `json:"buzzwords" dhall:"buzzwords"`
	Location            Location `json:"location" dhall:"location"`
	Name                string   `json:"name" dhall:"name"`
	NotablePublications []Link   `json:"notablePublications" dhall:"notablePublications"`
	Tagline             string   `json:"tagline" dhall:"tagline"`
}
