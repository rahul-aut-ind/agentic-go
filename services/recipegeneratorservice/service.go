package recipegeneratorservice

import (
	"context"
	"fmt"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
	"github.com/firebase/genkit/go/plugins/googlegenai"
	"github.com/rahul-aut-ind/genkit-go/domain/models"
	"github.com/rahul-aut-ind/genkit-go/pkg/config"
)

const (
	RecipeFlowName = "recipeGeneratorFlow"
)

type (
	Client struct {
		genkit *genkit.Genkit
	}
)

func NewClient(cfg *config.Config) *Client {
	return &Client{
		genkit: genkit.Init(context.Background(),
			genkit.WithPlugins(&googlegenai.GoogleAI{
				APIKey: cfg.APIKey,
			}),
			genkit.WithDefaultModel(cfg.ModelName),
		),
	}
}

func (c *Client) GenerateRecipe(ctx context.Context, ingredient string, dietaryRestrictions string) (*models.RecipeOutput, error) {
	// Define a type-safe flow
	recipeFlow := genkit.DefineFlow(c.genkit, RecipeFlowName, c.recipeFlow)

	// Run the flow
	recipe, err := recipeFlow.Run(ctx, &models.RecipeInput{
		Ingredient:          ingredient,
		DietaryRestrictions: dietaryRestrictions,
	})
	if err != nil {
		return nil, fmt.Errorf("could not generate recipe: %v", err)
	}

	return recipe, nil
}

func (c *Client) recipeFlow(ctx context.Context, input *models.RecipeInput) (*models.RecipeOutput, error) {
	dietaryRestrictions := input.DietaryRestrictions
	if dietaryRestrictions == "" {
		dietaryRestrictions = "none"
	}

	prompt := fmt.Sprintf(models.RecipePrompt, input.Ingredient, dietaryRestrictions)

	// Generate structured data with type safety
	recipe, _, err := genkit.GenerateData[models.RecipeOutput](ctx,
		c.genkit,
		ai.WithPrompt(prompt),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to generate recipe: %w", err)
	}

	return recipe, nil
}
