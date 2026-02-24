package postgres

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

// NewConnection: Postgres bağlantısını açar ve ping atarak kontrol eder
func NewConnection(dbUrl string) *sql.DB {
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatalf("failed to open db connection: %v", err)
	}

	// DB gerçekten ayakta mı? Test ediyoruz.
	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping db: %v", err)
	}

	log.Println("✅ Connected to Postgres")
	return db
}
