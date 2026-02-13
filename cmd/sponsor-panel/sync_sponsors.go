package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// graphqlSponsorsResponse represents the GraphQL response for sponsorshipsAsMaintainer.
type graphqlSponsorsResponse struct {
	Data struct {
		User struct {
			SponsorshipsAsMaintainer struct {
				PageInfo struct {
					HasNextPage bool   `json:"hasNextPage"`
					EndCursor   string `json:"endCursor"`
				} `json:"pageInfo"`
				Nodes []struct {
					SponsorEntity struct {
						Typename string `json:"__typename"`
						Login    string `json:"login"`
					} `json:"sponsorEntity"`
					Tier struct {
						Name                string `json:"name"`
						MonthlyPriceInCents int    `json:"monthlyPriceInCents"`
					} `json:"tier"`
					IsActive bool `json:"isActive"`
				} `json:"nodes"`
			} `json:"sponsorshipsAsMaintainer"`
		} `json:"user"`
	} `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

// syncSponsors performs a single sync of all sponsors from GitHub.
func syncSponsors(ctx context.Context, pool *pgxpool.Pool, ghToken string) error {
	slog.Info("syncSponsors: starting sponsor sync")

	allSponsors := make([]string, 0)

	after := ""
	for {
		// Build GraphQL query
		query := fmt.Sprintf(`query {
			user(login: "Xe") {
				sponsorshipsAsMaintainer(first: 100, after: %s, activeOnly: true) {
					pageInfo { hasNextPage, endCursor }
					nodes {
						sponsorEntity {
							__typename
							... on User { login }
							... on Organization { login }
						}
						tier { name, monthlyPriceInCents }
						isActive
					}
				}
			}
		}`, formatGraphQLString(after))

		reqBody := map[string]any{"query": query}
		bodyBytes, err := json.Marshal(reqBody)
		if err != nil {
			slog.Error("syncSponsors: failed to marshal request", "err", err)
			return err
		}

		req, err := http.NewRequestWithContext(ctx, "POST", "https://api.github.com/graphql", strings.NewReader(string(bodyBytes)))
		if err != nil {
			slog.Error("syncSponsors: failed to create request", "err", err)
			return err
		}

		req.Header.Set("Authorization", "Bearer "+ghToken)
		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			slog.Error("syncSponsors: request failed", "err", err)
			return err
		}

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			slog.Error("syncSponsors: GraphQL API error", "status", resp.StatusCode, "body", string(body))
			return fmt.Errorf("GraphQL API returned status %d: %s", resp.StatusCode, string(body))
		}

		var result graphqlSponsorsResponse
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			resp.Body.Close()
			slog.Error("syncSponsors: failed to decode response", "err", err)
			return err
		}
		resp.Body.Close()

		// Check for GraphQL errors
		if len(result.Errors) > 0 {
			slog.Error("syncSponsors: GraphQL errors", "errors", result.Errors)
			return fmt.Errorf("GraphQL error: %s", result.Errors[0].Message)
		}

		// Process sponsors
		for _, node := range result.Data.User.SponsorshipsAsMaintainer.Nodes {
			login := node.SponsorEntity.Login
			if login == "" {
				continue
			}

			sponsor := &SponsorUsername{
				Username:           login,
				EntityType:         node.SponsorEntity.Typename,
				MonthlyAmountCents: node.Tier.MonthlyPriceInCents,
				TierName:           node.Tier.Name,
				IsActive:           true,
			}

			// Upsert to database
			if err := upsertSponsorUsername(ctx, pool, sponsor); err != nil {
				slog.Error("syncSponsors: failed to upsert sponsor", "err", err, "username", login)
				continue
			}

			allSponsors = append(allSponsors, login)
		}

		// Check for next page
		if !result.Data.User.SponsorshipsAsMaintainer.PageInfo.HasNextPage {
			break
		}
		after = result.Data.User.SponsorshipsAsMaintainer.PageInfo.EndCursor
	}

	// Mark inactive sponsors not in current fetch
	inactiveCount, err := markInactiveSponsorsNotIn(ctx, pool, allSponsors)
	if err != nil {
		slog.Error("syncSponsors: failed to mark inactive sponsors", "err", err)
		return err
	}

	slog.Info("syncSponsors: sync completed",
		"active_sponsors", len(allSponsors),
		"marked_inactive", inactiveCount)

	return nil
}

// formatGraphQLString formats a string for GraphQL (with quotes, or null if empty).
func formatGraphQLString(s string) string {
	if s == "" {
		return "null"
	}
	return fmt.Sprintf(`"%s"`, s)
}

// startSyncLoop runs the sync immediately, then every hour.
func startSyncLoop(ctx context.Context, pool *pgxpool.Pool, ghToken string) {
	// Run initial sync immediately
	slog.Info("startSyncLoop: running initial sponsor sync")
	if err := syncSponsors(ctx, pool, ghToken); err != nil {
		slog.Error("startSyncLoop: initial sync failed", "err", err)
	}

	// Run every hour
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			slog.Info("startSyncLoop: stopping sync loop", "reason", ctx.Err())
			return
		case <-ticker.C:
			slog.Info("startSyncLoop: running scheduled sponsor sync")
			if err := syncSponsors(ctx, pool, ghToken); err != nil {
				slog.Error("startSyncLoop: scheduled sync failed", "err", err)
			}
		}
	}
}
