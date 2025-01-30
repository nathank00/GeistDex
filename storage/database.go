package storage

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() {
	var err error
	DB, err = sql.Open("sqlite3", "geistdex.db")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	// Create table for liquidity tracking
	query := `
	CREATE TABLE IF NOT EXISTS liquidity (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		pair TEXT NOT NULL,
		reserve0 TEXT NOT NULL,
		reserve1 TEXT NOT NULL,
		timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
	);`
	_, err = DB.Exec(query)
	if err != nil {
		log.Fatalf("Failed to create liquidity table: %v", err)
	}
}
