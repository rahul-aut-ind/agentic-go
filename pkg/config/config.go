package config

type (
	Config struct {
		ModelName string
		APIKey    string
	}
)

func NewConfig() *Config {
	// read from env
	return &Config{
		ModelName: "googleai/gemini-2.5-flash",
		APIKey:    "AIzaSyDYTT7p0re9lX1rBW4Nu_JZdwgZ5bq3M0w",
	}
}
