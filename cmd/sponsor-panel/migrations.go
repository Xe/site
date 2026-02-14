package main

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

const migrationSchema = `
-- Users table: GitHub accounts + sponsorship data
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    github_id BIGINT UNIQUE NOT NULL,
    login TEXT NOT NULL UNIQUE,
    avatar_url TEXT,
    name TEXT,
    email TEXT,

    -- Sponsorship data from GraphQL (cached)
    sponsorship_data JSONB,
    last_sponsorship_check TIMESTAMP DEFAULT NOW(),

    -- Timestamps
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Logo submissions: Simple tracking only
CREATE TABLE IF NOT EXISTS logo_submissions (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,

    company_name TEXT NOT NULL,
    website TEXT NOT NULL,
    logo_url TEXT,
    github_issue_url TEXT,
    github_issue_number INTEGER,

    submitted_at TIMESTAMP DEFAULT NOW()
);

-- Indexes for common queries
CREATE INDEX IF NOT EXISTS idx_users_github_id ON users(github_id);
CREATE INDEX IF NOT EXISTS idx_users_login ON users(login);
CREATE INDEX IF NOT EXISTS idx_logo_user_id ON logo_submissions(user_id);

-- GitHub sponsor usernames: synced list of all sponsors (users + orgs)
CREATE TABLE IF NOT EXISTS github_sponsor_usernames (
    id SERIAL PRIMARY KEY,
    username TEXT NOT NULL UNIQUE,           -- GitHub login (user or org)
    entity_type TEXT NOT NULL,               -- 'User' or 'Organization'
    monthly_amount_cents INTEGER DEFAULT 0,  -- Sponsorship tier amount
    tier_name TEXT,                          -- Tier name for display
    is_active BOOLEAN DEFAULT TRUE,          -- Active sponsorship flag
    synced_at TIMESTAMP DEFAULT NOW(),       -- Last sync timestamp
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_sponsor_usernames ON github_sponsor_usernames(username);
CREATE INDEX IF NOT EXISTS idx_sponsor_active ON github_sponsor_usernames(is_active);
`

// runMigrations executes the database schema migration.
func runMigrations(ctx context.Context, pool *pgxpool.Pool) error {
	slog.Info("running database migrations")
	_, err := pool.Exec(ctx, migrationSchema)
	if err != nil {
		return err
	}
	slog.Info("database migrations completed")
	return nil
}
