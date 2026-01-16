package models

type (
	IterativeRefinementRequest struct {
		Topic string `json:"topic"`
	}

	Evaluation struct {
		Critique  string `json:"critique"`
		Satisfied bool   `json:"satisfied"`
	}
)

const (
	GenerateBlogPrompt = `Write a short, single-paragraph (100-150 words) blog post about: %v.`
	CritiquePrompt     = `Critique the following blog post. Is it clear, concise, and engaging? Provide specific feedback for improvement. Post: "%v"`
	RefinePrompt       = `Revise the following blog post based on the feedback provided.\n Post: "%v" \n Feedback: "%v"`
)
