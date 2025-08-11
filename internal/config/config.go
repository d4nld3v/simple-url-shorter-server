package config

import "os"

type Config struct {
	Addr     string
	DBDriver string
	DBSource string
}

func Load() Config {
	return Config{
		Addr:     getEnv("ADDR", ":8080"),
		DBDriver: "sqlite3",
		DBSource: getEnv("DB_SOURCE", "url_shortener.db"),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
