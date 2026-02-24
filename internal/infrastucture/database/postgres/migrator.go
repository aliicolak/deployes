package postgres

import (
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigrations(databaseURL string) {
	log.Println("🔄 Running database migrations...")

	// Path relative to project root
	m, err := migrate.New(
		"file://internal/infrastucture/database/migrations",
		databaseURL,
	)
	if err != nil {
		log.Fatalf("❌ Migration initialization failed: %v", err)
	}

	if err := m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			log.Println("✅ No new migrations to apply")
		} else {
			log.Fatalf("❌ Migration failed: %v", err)
		}
	} else {
		log.Println("✅ Database migrations applied successfully")
	}
}
