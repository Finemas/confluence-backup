// config.go
package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	BaseURL  string
	Email    string // optional for Bearer
	Token    string
	SpaceKey string
}

func LoadConfig() Config {
	LoadEnv()
	return Config{
		BaseURL:  GetEnv("DRMAX_CONFLUENCE_BASE"),
		Email:    GetEnv("DRMAX_CONFLUENCE_EMAIL"),
		Token:    GetEnv("DRMAX_CONFLUENCE_TOKEN"),
		SpaceKey: GetEnv("DRMAX_SPACE_KEY"),
	}
}

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func GetEnv(key string) string {
	return os.Getenv(key)
}
