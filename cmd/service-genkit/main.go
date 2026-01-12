package main

import (
	"context"
	"fmt"
	"time"

	"github.com/rahul-aut-ind/genkit-go/pkg/config"
	"github.com/rahul-aut-ind/genkit-go/pkg/logger"
	"github.com/rahul-aut-ind/genkit-go/services/recipegeneratorservice"
)

var (
	ingredient          = "avocado"
	dietaryRestrictions = "vegetarian"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// initialize everything
	env := config.NewConfig()
	log := logger.New()
	service := recipegeneratorservice.NewClient(ctx, env)

	// generate recipe
	resp, err := service.GenerateRecipe(ctx, ingredient, dietaryRestrictions)
	if err != nil {
		fmt.Println("Error generating recipe:", err)
		return
	}

	// Print the structured recipe
	log.PrettifyJSON(resp)

}
