package config

import "os"

type Config struct {
	Addr       string
	DBDriver   string
	DBSource   string
	RateLimit  int
	BurstLimit int
}

func Load() Config {
	return Config{
		Addr:       getEnv("ADDR", ":8080"),
		DBDriver:   "sqlite3",
		DBSource:   getEnv("DB_SOURCE", "url_shortener.db"),
		RateLimit:  10, // 10 requests per second
		BurstLimit: 5,  // 5 requests burst capacity
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
