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
	dietaryRestrictions = "none"
)

func main() {
	// initialize everything
	env := config.NewConfig()
	log := logger.New()
	service := recipegeneratorservice.NewClient(context.Background(), env)

	fmt.Println("Enter ingredient:")
	fmt.Scanf("%s", &ingredient)
	fmt.Println("Enter dietary restrictions if any (leave blank for none):")
	fmt.Scanf("%s", &dietaryRestrictions)
	fmt.Println("ok, generating a recipe for you with ingredient ", ingredient, " and dietary restrictions ", dietaryRestrictions)
	// generate recipe
	start := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	resp, err := service.GenerateRecipe(ctx, ingredient, dietaryRestrictions)
	if err != nil {
		fmt.Println("Error generating recipe:", err)
		return
	}

	// Print the structured recipe
	log.PrettifyJSON(resp)
	log.Infof("Recipe generated in %s", time.Since(start))

}
