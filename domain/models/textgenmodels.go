package models

type (
	RouterRequest struct {
		Query string `json:"query"`
	}

	Intent struct {
		Intent string `json:"intent" jsonschema_enum:"question,creative"`
	}
)

const (
	IntentGenPrompt      = `Classify the user's query as either a 'question' or a 'creative' request. Query: %v`
	QuestionIntentPrompt = `Answer the following question: %v`
	CreativeIntentPrompt = `Write a short poem about: %v`
)
