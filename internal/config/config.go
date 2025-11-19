package config

import "os"

type Config struct {
	App      app
	Database database
	Cache    cache
	Bible    bible
}

func NewConfig() *Config {
	return &Config{
		App:      newApp(),
		Database: newDatabase(),
		Cache:    newCache(),
		Bible:    newBible(),
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	return value
}
