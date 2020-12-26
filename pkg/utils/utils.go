package utils

import (
	"os"

	"github.com/joho/godotenv"
)

func LoadConfig() error {
	err := LoadEnv()
	if err != nil {
		return err
	}

	return nil
}

func LoadEnv() (err error) {
	p := os.Getenv("ENV")
	if p == "" {
		p = ".env"
	}

	return godotenv.Load(p)
}
