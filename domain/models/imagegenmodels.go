package models

type (
	ImageGeneratorRequest struct {
		Concept string `json:"concept"`
	}
	ImageGeneratorResponse struct {
		ImageURL string `json:"image_url"`
	}
)

const (
	ImageConceptPrompt = `Create a detailed, artistic prompt for an image generation model. The concept is: %s`
)
