package bloggeneratorservice

import (
	"context"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
	"github.com/firebase/genkit/go/plugins/googlegenai"
	"github.com/rahul-aut-ind/genkit-go/domain/models"
	"github.com/rahul-aut-ind/genkit-go/pkg/config"
	"github.com/rahul-aut-ind/genkit-go/pkg/logger"
)

const (
	IterativeRefinementFlowName = "iterativeRefinementFlow"
)

type (
	Client struct {
		genkit       *genkit.Genkit
		generalModel string
		log          *logger.Logger
	}
)

const (
	NoOfIterations = 3
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

func (c *Client) GenerateBlog(ctx context.Context, topic string) (string, error) {
	iterativeRefinementFlow := genkit.DefineFlow(c.genkit, IterativeRefinementFlowName, c.iterativeRefinementFlow)
	resp, err := iterativeRefinementFlow.Run(ctx, &models.IterativeRefinementRequest{
		Topic: topic,
	})
	if err != nil {
		return "", err
	}
	return resp, nil
}

func (c *Client) iterativeRefinementFlow(ctx context.Context, req *models.IterativeRefinementRequest) (string, error) {

	// Step 1: Generate the initial draft.
	resp, err := genkit.Generate(ctx, c.genkit,
		ai.WithPrompt(models.GenerateBlogPrompt, req.Topic),
	)
	if err != nil {
		return "", err
	}
	content := resp.Text()
	c.log.PrettifyJSON("Initial Content: " + content)

	// Step 2: Iteratively refine the content.
	for range NoOfIterations {
		// The "Evaluator" provides feedback.
		eval, _, err := genkit.GenerateData[models.Evaluation](ctx, c.genkit,
			ai.WithPrompt(models.CritiquePrompt, content),
		)
		if err != nil {
			return "", err
		}
		c.log.PrettifyJSON(eval)
		if eval.Satisfied {
			break
		}

		// The "Optimizer" refines the content based on feedback.
		resp, err := genkit.Generate(ctx, c.genkit,
			ai.WithPrompt(models.RefinePrompt, content, eval.Critique),
		)
		if err != nil {
			return "", err
		}
		content = resp.Text()
		c.log.PrettifyJSON("Refined Content: " + content)
	}
	return content, nil
}
