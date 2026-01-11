package config

import (
	"os"

	"github.com/joho/godotenv"
)

type (
	Config struct {
		ModelName string
		APIKey    string
	}
)

func NewConfig() *Config {
	path, _ := os.Getwd()
	err := godotenv.Load(path + "/.env")
	if err != nil {
		return &Config{
			ModelName: "googleai/gemini-2.5-flash",
			APIKey:    "generate from aistudio",
		}
	}
	return &Config{
		ModelName: os.Getenv("MODEL_NAME"),
		APIKey:    os.Getenv("API_KEY"),
	}
}
