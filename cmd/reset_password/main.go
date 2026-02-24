package main

import (
	"database/sql"
	"deployes/pkg/utils"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping DB: %v", err)
	}

	password := "password"
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		log.Fatal(err)
	}

	email := "admin@deployes.com"
	var exists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email=$1)", email).Scan(&exists)
	if err != nil {
		log.Fatalf("Failed to check if user exists: %v", err)
	}

	if exists {
		_, err = db.Exec("UPDATE users SET password=$1 WHERE email=$2", hashedPassword, email)
		if err != nil {
			log.Fatalf("Failed to update password: %v", err)
		}
		fmt.Printf("Password for %s updated successfully to '%s'\n", email, password)
	} else {
		log.Printf("User %s not found. Creating...", email)
		id := uuid.NewString()
		now := time.Now()
		_, err = db.Exec(`
			INSERT INTO users (id, email, password, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5)
		`, id, email, hashedPassword, now, now)
		if err != nil {
			log.Fatalf("Failed to create user: %v", err)
		}
		fmt.Printf("User %s created successfully with password '%s'\n", email, password)
	}
}
