package utils

import (
	"os"

	"github.com/joho/godotenv"
)

func GetDotEnv(key string) string {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	value := os.Getenv(key)
	return value
}