package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"sort"
	"time"

	"github.com/mxpv/patreon-go"
	"golang.org/x/oauth2"
	"within.website/ln"
)

func NewPatreonClient() (*patreon.Client, error) {
	for _, name := range []string{"CLIENT_ID", "CLIENT_SECRET", "ACCESS_TOKEN", "REFRESH_TOKEN"} {
		if os.Getenv("PATREON_"+name) == "" {
			return nil, fmt.Errorf("wanted envvar PATREON_%s", name)
		}
	}

	config := oauth2.Config{
		ClientID:     os.Getenv("PATREON_CLIENT_ID"),
		ClientSecret: os.Getenv("PATREON_CLIENT_SECRET"),
		Endpoint: oauth2.Endpoint{
			AuthURL:  patreon.AuthorizationURL,
			TokenURL: patreon.AccessTokenURL,
		},
		Scopes: []string{"users", "campaigns", "pledges", "pledges-to-me", "my-campaign"},
	}

	token := oauth2.Token{
		AccessToken:  os.Getenv("PATREON_ACCESS_TOKEN"),
		RefreshToken: os.Getenv("PATREON_REFRESH_TOKEN"),
		// Must be non-nil, otherwise token will not be expired
		Expiry: time.Now().Add(90 * 24 * time.Hour),
	}

	tc := config.Client(context.Background(), &token)

	trans := tc.Transport
	tc.Transport = lnLoggingTransport{next: trans}
	client := patreon.NewClient(tc)

	return client, nil
}

func GetPledges(pc *patreon.Client) ([]string, error) {
	campaign, err := pc.FetchCampaign()
	if err != nil {
		return nil, fmt.Errorf("campaign fetch error: %w", err)
	}

	campaignID := campaign.Data[0].ID

	cursor := ""
	var result []string

	for {
		pledgesResponse, err := pc.FetchPledges(campaignID, patreon.WithPageSize(25), patreon.WithCursor(cursor))
		if err != nil {
			return nil, err
		}

		users := make(map[string]*patreon.User)
		for _, item := range pledgesResponse.Included.Items {
			u, ok := item.(*patreon.User)
			if !ok {
				continue
			}

			users[u.ID] = u
		}

		for _, pledge := range pledgesResponse.Data {
			pid := pledge.Relationships.Patron.Data.ID
			patronFullName := users[pid].Attributes.FullName

			result = append(result, patronFullName)
		}

		cursor = pledgesResponse.Links.Next
		if cursor == "" {
			break
		}
	}

	sort.Strings(result)
	return result, nil
}

type lnLoggingTransport struct{ next http.RoundTripper }

func (l lnLoggingTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	ctx := r.Context()
	f := ln.F{
		"url":       r.URL.String(),
		"has_token": r.Header.Get("Authorization") != "",
	}

	resp, err := l.next.RoundTrip(r)
	if err != nil {
		return nil, err
	}

	f["status"] = resp.Status

	ln.Log(ctx, f)

	return resp, nil
}
