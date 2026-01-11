package main

import (
	"context"
	"fmt"

	"github.com/rahul-aut-ind/genkit-go/pkg/config"
	"github.com/rahul-aut-ind/genkit-go/pkg/logger"
	"github.com/rahul-aut-ind/genkit-go/services/recipegeneratorservice"
)

var (
	ingredient          = "avocado"
	dietaryRestrictions = "vegetarian"
)

func main() {
	ctx := context.Background()

	// initialize config
	env := config.NewConfig()
	logger := logger.New()

	// Initialize service
	service := recipegeneratorservice.NewClient(env)

	// generate recipe
	resp, err := service.GenerateRecipe(ctx, ingredient, dietaryRestrictions)
	if err != nil {
		fmt.Println("Error generating recipe:", err)
		return
	}

	// Print the structured recipe
	logger.PrettifyJSON(resp)

}
