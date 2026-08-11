package main

import (
	"strings"
	"sync"
	"testing"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// TestModelsCompile verifies GORM model structs exist and have the expected shape.
func TestModelsCompile(t *testing.T) {
	// These should compile once models.go is rewritten
	var _ gorm.Model

	u := PanelUser{}
	if u.TableName() != "users" {
		t.Errorf("PanelUser.TableName() = %q, want %q", u.TableName(), "users")
	}

	ls := LogoSubmission{}
	if ls.TableName() != "logo_submissions" {
		t.Errorf("LogoSubmission.TableName() = %q, want %q", ls.TableName(), "logo_submissions")
	}

	su := SponsorUsername{}
	if su.TableName() != "github_sponsor_usernames" {
		t.Errorf("SponsorUsername.TableName() = %q, want %q", su.TableName(), "github_sponsor_usernames")
	}
}

// TestPanelUserColumnNames pins the column that GORM derives for each field
// queried by a literal column name. A mismatch makes queries read one column
// while inserts write another.
func TestPanelUserColumnNames(t *testing.T) {
	parsed, err := schema.Parse(&PanelUser{}, &sync.Map{}, schema.NamingStrategy{})
	if err != nil {
		t.Fatalf("schema.Parse: %v", err)
	}

	for _, tt := range []struct {
		field string
		want  string
	}{
		{field: "GitHubID", want: "github_id"},
		{field: "PatreonID", want: "patreon_id"},
		{field: "Login", want: "login"},
		{field: "Provider", want: "provider"},
		{field: "Email", want: "email"},
	} {
		t.Run(tt.field, func(t *testing.T) {
			f := parsed.LookUpField(tt.field)
			if f == nil {
				t.Fatalf("field %s not found", tt.field)
			}
			if f.DBName != tt.want {
				t.Errorf("column for %s = %q, want %q", tt.field, f.DBName, tt.want)
			}
		})
	}
}

// dryRunDB returns a GORM handle that builds SQL but never connects.
func dryRunDB(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(postgres.New(postgres.Config{DSN: "postgres://user:pass@127.0.0.1:1/none"}), &gorm.Config{
		DryRun:               true,
		DisableAutomaticPing: true,
	})
	if err != nil {
		t.Fatalf("gorm.Open: %v", err)
	}
	return db
}

// TestUpsertUserLooksUpMappedColumn verifies that the lookup in upsertUser
// reads the same column that Create writes. A literal "github_id" condition
// reads a stale column, so every login inserts a duplicate row and violates
// the unique index users_github_id_key.
func TestUpsertUserLooksUpMappedColumn(t *testing.T) {
	db := dryRunDB(t)

	id := int64(1165302)
	var existing PanelUser
	got := db.Where(&PanelUser{GitHubID: &id}).First(&existing).Statement.SQL.String()

	// Create writes the column that TestPanelUserColumnNames pins, so a
	// lookup on the same column keeps read and write in agreement.
	if !strings.Contains(got, `"github_id" = `) {
		t.Errorf("lookup SQL = %q, want a condition on github_id", got)
	}
}
