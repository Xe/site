package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/facebookgo/flagenv"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/joho/godotenv/autoload"
)

var databaseURL = flag.String("database-url", "", "Database URL")

func main() {
	flagenv.Parse()
	flag.Parse()

	if *databaseURL == "" {
		fmt.Fprintln(os.Stderr, "database-url is required")
		os.Exit(1)
	}

	db, err := sql.Open("pgx", *databaseURL)
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping: %v", err)
	}

	statements := []string{
		// Users table: drop constraints (unique indexes backed by constraints)
		`ALTER TABLE users DROP CONSTRAINT IF EXISTS users_github_id_key`,
		`ALTER TABLE users DROP CONSTRAINT IF EXISTS users_patreon_id_key`,
		// Users table: drop plain indexes
		`DROP INDEX IF EXISTS idx_users_github_id`,
		`DROP INDEX IF EXISTS idx_users_login`,
		`DROP INDEX IF EXISTS idx_users_patreon_id`,
		`DROP INDEX IF EXISTS idx_users_provider_login`,
		// Sponsor usernames table: drop constraint
		`ALTER TABLE github_sponsor_usernames DROP CONSTRAINT IF EXISTS github_sponsor_usernames_username_key`,
		// Sponsor usernames table: drop plain indexes
		`DROP INDEX IF EXISTS idx_sponsor_active`,
		`DROP INDEX IF EXISTS idx_sponsor_usernames`,
	}

	for _, stmt := range statements {
		fmt.Printf("  %s\n", stmt)
		if _, err := db.Exec(stmt); err != nil {
			log.Fatalf("failed: %v", err)
		}
	}

	fmt.Println("done — GORM AutoMigrate will recreate indexes on next startup")
}
