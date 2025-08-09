package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDatabase() error {
	var err error
	DB, err = sql.Open("sqlite3", "./url_shortener.db")
	if err != nil {
		return err
	}

	// Create table if not exists
	createTableSQL := `CREATE TABLE IF NOT EXISTS urls (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		original_url TEXT NOT NULL,
		shorten_id TEXT NOT NULL UNIQUE
	);`
	_, err = DB.Exec(createTableSQL)
	if err != nil {
		return err
	}

	return nil
}

func CloseDatabase() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}
