package models

import "context"

type (
	RecipeGenerator interface {
		GenerateRecipe(ctx context.Context, ingredient string, dietaryRestrictions string) (*RecipeOutput, error)
	}

	RecipeInput struct {
		Ingredient          string `json:"ingredient" jsonschema:"description=Main ingredient or cuisine type"`
		DietaryRestrictions string `json:"dietaryRestrictions,omitempty" jsonschema:"description=Any dietary restrictions"`
	}

	RecipeOutput struct {
		Title        string   `json:"title"`
		Description  string   `json:"description"`
		PrepTime     string   `json:"prepTime"`
		CookTime     string   `json:"cookTime"`
		Servings     int      `json:"servings"`
		Ingredients  []string `json:"ingredients"`
		Instructions []string `json:"instructions"`
		Tips         []string `json:"tips,omitempty"`
	}
)

const (
	RecipePrompt = `Create a recipe with the following requirements:
		Main ingredient: %s
		Dietary restrictions: %s`
)
