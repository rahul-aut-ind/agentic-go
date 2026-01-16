package models

type (
	MarketingIdeaResponse struct {
		Name    string `json:"name"`
		Tagline string `json:"tagline"`
	}

	MarketingIdeaRequest struct {
		ProductName string `json:"productName"`
		Location    string `json:"location" jsonschema_description:"The city name, e.g. San Francisco"`
	}

	WeatherToolInput struct {
		Location string `json:"location" jsonschema_description:"The city name, e.g. San Francisco"`
	}
)

const (
	CreativeNamePrompt    = `Generate a creative name for a new product: %s. Product is to be launched in %s where weather is : %s.`
	CreativeTaglinePrompt = `Generate a catchy tagline for a new product: %s. Product is to be launched in %s where weather is : %s.`
	GetWeatherPrompt      = `Give me the weather in the location: %s.`
)
