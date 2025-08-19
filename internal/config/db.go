package config

import (
	"database/sql"
	"fmt"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

// DB is a singleton instance of the database connection
var (
	DB   *sql.DB
	once sync.Once
)

func InitDB() error {
	var err error
	once.Do(func() {
		DB, err = sql.Open("sqlite3", "./url_shortener.db")
		if err != nil {
			return
		}

		DB.SetMaxOpenConns(25)
		DB.SetMaxIdleConns(5)

		if pingErr := DB.Ping(); pingErr != nil {
			err = pingErr
			return
		}

		createTableSQL := `CREATE TABLE IF NOT EXISTS urls (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            original_url TEXT NOT NULL,
            shorten_id TEXT NOT NULL UNIQUE,
            clicks INTEGER DEFAULT 0,
            created_at DATETIME DEFAULT CURRENT_TIMESTAMP
        );`
		_, err = DB.Exec(createTableSQL)

		if err != nil {
			fmt.Println("Error creating table:", err)
			return
		}

		fmt.Println("Database initialized successfully")
	})

	return err
}

func GetDB() (*sql.DB, error) {

	error := InitDB()
	if error != nil {
		fmt.Println("Error initializing database:", error)
		return nil, error
	}

	if DB == nil {
		return nil, fmt.Errorf("database not initialized")
	}
	return DB, nil
}

func CloseDatabase() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}
