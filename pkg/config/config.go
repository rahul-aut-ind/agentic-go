package config

import (
	"os"

	"github.com/joho/godotenv"
)

type (
	Config struct {
		GeneralModelName  string
		ImageGenModelName string
		APIKey            string
	}
)

func NewConfig() *Config {
	path, _ := os.Getwd()
	err := godotenv.Load(path + "/.env")
	if err != nil {
		return &Config{
			GeneralModelName:  "googleai/gemini-2.5-flash",
			ImageGenModelName: "googleai/imagen-3.0-generate-002",
			APIKey:            "generate from aistudio",
		}
	}
	return &Config{
		GeneralModelName:  os.Getenv("GENERAL_MODEL_NAME"),
		ImageGenModelName: os.Getenv("IMAGE_GEN_MODEL_NAME"),
		APIKey:            os.Getenv("API_KEY"),
	}
}
