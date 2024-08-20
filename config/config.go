package config

import (
	"github.com/joho/godotenv"
)

func Load(path string) {
	err := godotenv.Load(path)
	if err != nil {
		panic(err)
	}
}
