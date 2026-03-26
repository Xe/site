package main

import (
	"testing"

	"gorm.io/gorm"
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
