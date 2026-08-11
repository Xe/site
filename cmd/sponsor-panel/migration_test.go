package main

import (
	"os"
	"strings"
	"testing"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// productionSchema recreates the users table as it stood before the GitHub ID
// columns were merged: a legacy github_id column plus the git_hub_id column
// that GORM derived from the field name.
const productionSchema = `
DROP TABLE IF EXISTS logo_submissions;
DROP TABLE IF EXISTS users;

CREATE TABLE users (
    id                     serial PRIMARY KEY,
    github_id              bigint,
    login                  text NOT NULL,
    avatar_url             text,
    name                   text,
    email                  text,
    sponsorship_data       jsonb,
    last_sponsorship_check timestamp,
    created_at             timestamp,
    updated_at             timestamp,
    patreon_id             text,
    provider               text NOT NULL DEFAULT 'github',
    git_hub_id             bigint,
    thoth_user_id          text
);

CREATE UNIQUE INDEX idx_users_provider_login ON users (provider, login);
CREATE UNIQUE INDEX users_github_id_key ON users (git_hub_id);
CREATE UNIQUE INDEX users_patreon_id_key ON users (patreon_id);

INSERT INTO users (github_id, git_hub_id, patreon_id, login, provider) VALUES
    (1111, NULL, NULL, 'legacy-only', 'github'),
    (NULL, 2222, NULL, 'gorm-only',   'github'),
    (NULL, NULL,  'p1', 'patreon-only', 'patreon');
`

// testDB connects to a throwaway PostgreSQL instance. Set SPONSOR_PANEL_TEST_DSN
// to run the migration tests.
func testDB(t *testing.T) *gorm.DB {
	t.Helper()

	dsn := os.Getenv("SPONSOR_PANEL_TEST_DSN")
	if dsn == "" {
		t.Skip("SPONSOR_PANEL_TEST_DSN is not set")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: logger.Discard})
	if err != nil {
		t.Fatalf("gorm.Open: %v", err)
	}

	if err := db.Exec(productionSchema).Error; err != nil {
		t.Fatalf("cannot create the test schema: %v", err)
	}

	return db
}

// columnNames returns the columns of the users table.
func columnNames(t *testing.T, db *gorm.DB) map[string]bool {
	t.Helper()

	var names []string
	err := db.Raw(`SELECT column_name FROM information_schema.columns
	               WHERE table_name = 'users' AND table_schema = current_schema()`).Scan(&names).Error
	if err != nil {
		t.Fatalf("cannot read the users columns: %v", err)
	}

	out := map[string]bool{}
	for _, n := range names {
		out[n] = true
	}
	return out
}

// TestConsolidateGitHubIDColumn verifies the merge against a real database.
func TestConsolidateGitHubIDColumn(t *testing.T) {
	db := testDB(t)

	if err := consolidateGitHubIDColumn(db); err != nil {
		t.Fatalf("consolidateGitHubIDColumn: %v", err)
	}

	cols := columnNames(t, db)
	if cols["git_hub_id"] {
		t.Error("users.git_hub_id still exists after the merge")
	}
	if !cols["github_id"] {
		t.Fatal("users.github_id is missing after the merge")
	}

	for _, tt := range []struct {
		login string
		want  int64
	}{
		{login: "legacy-only", want: 1111},
		{login: "gorm-only", want: 2222},
	} {
		t.Run(tt.login, func(t *testing.T) {
			var got int64
			if err := db.Raw(`SELECT github_id FROM users WHERE login = ?`, tt.login).Scan(&got).Error; err != nil {
				t.Fatalf("cannot read github_id: %v", err)
			}
			if got != tt.want {
				t.Errorf("github_id for %s = %d, want %d", tt.login, got, tt.want)
			}
		})
	}

	// The unique index must follow the column, or duplicate accounts return.
	var indexed string
	err := db.Raw(`SELECT indexdef FROM pg_indexes
	               WHERE tablename = 'users' AND indexname = 'users_github_id_key'`).Scan(&indexed).Error
	if err != nil {
		t.Fatalf("cannot read the index definition: %v", err)
	}
	if indexed == "" {
		t.Fatal("users_github_id_key is missing after the merge")
	}
	if !strings.Contains(indexed, "(github_id)") {
		t.Errorf("users_github_id_key = %q, want an index on github_id", indexed)
	}

	// AutoMigrate must not re-create the derived column.
	if err := db.AutoMigrate(PanelModels()...); err != nil {
		t.Fatalf("AutoMigrate: %v", err)
	}
	if columnNames(t, db)["git_hub_id"] {
		t.Error("AutoMigrate re-created users.git_hub_id")
	}

	// A second run must do nothing.
	if err := consolidateGitHubIDColumn(db); err != nil {
		t.Fatalf("consolidateGitHubIDColumn is not idempotent: %v", err)
	}
}

// TestUpsertUserFindsMigratedRows verifies that a login after the merge updates
// the existing row instead of inserting a duplicate.
func TestUpsertUserFindsMigratedRows(t *testing.T) {
	db := testDB(t)

	if err := consolidateGitHubIDColumn(db); err != nil {
		t.Fatalf("consolidateGitHubIDColumn: %v", err)
	}
	if err := db.AutoMigrate(PanelModels()...); err != nil {
		t.Fatalf("AutoMigrate: %v", err)
	}

	var before int64
	db.Model(&PanelUser{}).Count(&before)

	for _, tt := range []struct {
		name  string
		id    int64
		login string
	}{
		{name: "row migrated from the legacy column", id: 1111, login: "legacy-only"},
		{name: "row written by GORM", id: 2222, login: "gorm-only"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			user := &PanelUser{
				GitHubID:        &tt.id,
				Login:           tt.login,
				Name:            "updated",
				SponsorshipData: `{"is_active":true,"monthly_amount_cents":5000}`,
			}
			if err := upsertUser(db, user); err != nil {
				t.Fatalf("upsertUser: %v", err)
			}
			if user.Name != "updated" {
				t.Errorf("Name = %q, want the update to survive", user.Name)
			}
		})
	}

	var after int64
	db.Model(&PanelUser{}).Count(&after)
	if after != before {
		t.Errorf("user count = %d, want %d; upsertUser inserted duplicates", after, before)
	}
}
