package main

import (
	"context"
	"fmt"
	"time"

	"github.com/rahul-aut-ind/genkit-go/pkg/config"
	"github.com/rahul-aut-ind/genkit-go/pkg/logger"
	"github.com/rahul-aut-ind/genkit-go/services/marketingtoolservice"
)

var (
	productName    = "Short Walks For Elderly"
	targetLocation = "Berlin"
)

func main() {
	// initialize everything
	env := config.NewConfig()
	log := logger.New()

	service := marketingtoolservice.NewClient(context.Background(), env, log)

	fmt.Println("Type your product name:")
	_, e := fmt.Scanf("%s", &productName)
	if e != nil {
		log.Warnf("Error reading product name: %w, taking default", e)
	}

	fmt.Println("Type your target location:")
	_, e = fmt.Scanf("%s", &targetLocation)
	if e != nil {
		log.Warnf("Error reading target location: %w, taking default", e)
	}
	fmt.Printf("\nok, generating a name and tagline for you based on the product name %s for launch in %s\n", productName, targetLocation)

	// generate text
	start := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	resp, err := service.GenerateMarketingIdea(ctx, productName, targetLocation)
	if err != nil {
		log.Errorf("Error generating response: %v", err)
		return
	}

	log.PrettifyJSON(resp)
	log.Infof("Response generated in %s", time.Since(start))
}
