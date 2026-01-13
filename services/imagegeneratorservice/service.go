package imagegeneratorservice

import (
	"context"
	"errors"
	"fmt"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
	"github.com/firebase/genkit/go/plugins/googlegenai"
	"github.com/rahul-aut-ind/genkit-go/domain/models"
	"github.com/rahul-aut-ind/genkit-go/pkg/config"
	"github.com/rahul-aut-ind/genkit-go/pkg/logger"
)

const (
	ImageGeneratorFlowName = "imageGeneratorFlow"
)

type (
	Client struct {
		genkit        *genkit.Genkit
		generalModel  string
		imageGenModel string
		log           *logger.Logger
	}
)

func NewClient(ctx context.Context, cfg *config.Config, log *logger.Logger) *Client {
	return &Client{
		genkit: genkit.Init(ctx,
			genkit.WithPlugins(&googlegenai.GoogleAI{
				APIKey: cfg.APIKey,
			}),
		),
		generalModel:  cfg.GeneralModelName,
		imageGenModel: cfg.ImageGenModelName,
		log:           log,
	}
}

func (c *Client) GenerateImage(ctx context.Context, concept string) (string, error) {
	// Define a type-safe flow
	recipeFlow := genkit.DefineFlow(c.genkit, ImageGeneratorFlowName, c.generateImageFlow)

	// Run the flow
	resp, err := recipeFlow.Run(ctx, &models.ImageGeneratorRequest{
		Concept: concept,
	})
	if err != nil {
		return "", fmt.Errorf("could not generate image: %v", err)
	}

	return resp, nil
}

func (c *Client) generateImageFlow(ctx context.Context, req *models.ImageGeneratorRequest) (string, error) {
	// Step 1: Use a text model to generate a rich image prompt based on concept
	promptResponse, err := genkit.Generate(ctx, c.genkit,
		ai.WithModelName(c.imageGenModel),
		ai.WithPrompt(fmt.Sprintf(models.ImageConceptPrompt, req.Concept)),
	)
	if err != nil {
		return "", err
	}
	imagePrompt := promptResponse.Text()
	// c.log.Infof("Generated image prompt: %s", imagePrompt)

	// Step 2: Use the generated prompt to create an image
	imageResponse, err := genkit.Generate(ctx, c.genkit,
		ai.WithModelName(c.generalModel),
		ai.WithPrompt(imagePrompt, nil),
	)
	if err != nil {
		return "", err
	}
	for _, m := range imageResponse.Message.Content {
		if m.IsMedia() {
			return m.Text, nil
		}
	}
	return "", errors.New("did not generate an image")
}
