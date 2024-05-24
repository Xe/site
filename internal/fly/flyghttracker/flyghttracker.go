package flyghttracker

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"strings"
	"time"

	"within.website/x/web"
)

var (
	flyghttrackerURL = flag.String("flyghttracker-url", "https://flyght-tracker.fly.dev/api/upcoming_events", "Flyghttracker URL")
)

// Date represents a date in the format "YYYY-MM-DD"
type Date struct {
	time.Time
}

// UnmarshalJSON parses a JSON string in the format "YYYY-MM-DD" to a Date
func (d *Date) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return err
	}
	d.Time = t
	return nil
}

// MarshalJSON returns a JSON string in the format "YYYY-MM-DD"
func (d Date) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.Time.Format("2006-01-02"))
}

// Event represents an event that members of DevRel will be attending.
type Event struct {
	Name      string   `json:"name"`
	URL       string   `json:"url"`
	StartDate Date     `json:"start_date"`
	EndDate   Date     `json:"end_date"`
	Location  string   `json:"location"`
	People    []string `json:"people"`
}

// Fetch new events from the Flyght Tracker URL.
//
// It returns a list of events that end in the future and that have "Xe" as one of the attendees.
func Fetch() ([]Event, error) {
	resp, err := http.Get(*flyghttrackerURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch flyghttracker events: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, web.NewError(http.StatusOK, resp)
	}

	var events []Event
	if err := json.NewDecoder(resp.Body).Decode(&events); err != nil {
		return nil, fmt.Errorf("failed to decode flyghttracker events: %w", err)
	}

	var result []Event

	for _, event := range events {
		if event.EndDate.Before(time.Now()) {
			continue
		}

		found := false
		for _, person := range event.People {
			if person == "Xe" {
				found = true
				break
			}
		}

		if found {
			result = append(result, event)
		}
	}

	return result, nil
}
