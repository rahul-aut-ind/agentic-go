package textgeneratorservice

import (
	"context"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
	"github.com/firebase/genkit/go/plugins/googlegenai"
	"github.com/rahul-aut-ind/genkit-go/domain/models"
	"github.com/rahul-aut-ind/genkit-go/pkg/config"
	"github.com/rahul-aut-ind/genkit-go/pkg/logger"
	"google.golang.org/genai"
)

const (
	ConditionalFlowName = "conditionalFlow"
)

type (
	Client struct {
		genkit       *genkit.Genkit
		generalModel string
		log          *logger.Logger
	}
)

func NewClient(ctx context.Context, cfg *config.Config, log *logger.Logger) *Client {
	return &Client{
		genkit: genkit.Init(ctx,
			genkit.WithPlugins(&googlegenai.GoogleAI{
				APIKey: cfg.APIKey,
			}),
			genkit.WithDefaultModel(cfg.GeneralModelName),
		),
		generalModel: cfg.GeneralModelName,
		log:          log,
	}
}

func (c *Client) GenerateText(ctx context.Context, query string) (string, error) {
	textGenFlow := genkit.DefineFlow(c.genkit, ConditionalFlowName, c.generateTextConditionalFlow)
	resp, err := textGenFlow.Run(ctx, &models.RouterRequest{
		Query: query,
	})
	if err != nil {
		return "", err
	}
	return resp, nil
}

func (c *Client) generateTextConditionalFlow(ctx context.Context, req *models.RouterRequest) (string, error) {
	// Step 1: Classify the user's intent
	intent, _, err := genkit.GenerateData[models.Intent](ctx, c.genkit,
		ai.WithPrompt(models.IntentGenPrompt, req.Query),
	)
	if err != nil {
		return "", err
	}

	c.log.Infof("Generated intent: %s", intent.Intent)

	// Step 2: Route based on the intent
	switch intent.Intent {
	case "question":
		// Handle as a straightforward question
		answerResponse, err := genkit.Generate(ctx, c.genkit,
			ai.WithPrompt(models.QuestionIntentPrompt, req.Query),
		)
		if err != nil {
			return "", err
		}
		return answerResponse.Text(), nil
	case "creative":
		// Handle as a creative writing prompt
		creativeResponse, err := genkit.Generate(ctx, c.genkit,
			ai.WithPrompt(models.CreativeIntentPrompt, req.Query),
			ai.WithModelName(c.generalModel), // or a different model
			ai.WithConfig(&genai.GenerateContentConfig{
				// MaxOutputTokens: 128,
				Temperature: genai.Ptr[float32](0.9),
				TopP:        genai.Ptr[float32](0.4),
				TopK:        genai.Ptr[float32](50.0),
			}),
		)
		if err != nil {
			return "", err
		}
		return creativeResponse.Text(), nil
	default:
		return "Sorry, I couldn't determine how to handle your request.", nil
	}
}
