package main

import (
	"context"
	"fmt"
	"time"

	"github.com/rahul-aut-ind/genkit-go/pkg/config"
	"github.com/rahul-aut-ind/genkit-go/pkg/logger"
	"github.com/rahul-aut-ind/genkit-go/services/textgeneratorservice"
)

var (
	query = "The sky is blue"
)

func main() {
	// initialize everything
	env := config.NewConfig()
	log := logger.New()

	service := textgeneratorservice.NewClient(context.Background(), env, log)

	fmt.Println("Type your query:")
	_, e := fmt.Scanf("%s", &query)
	if e != nil {
		log.Warnf("Error reading query: %w, taking default", e)
	}
	fmt.Println("ok, generating a response for you based on the query : ", query)

	// generate text
	start := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	resp, err := service.GenerateText(ctx, query)
	if err != nil {
		log.Errorf("Error generating response: %v", err)
		return
	}

	log.PrettifyJSON(resp)
	log.Infof("Response generated in %s", time.Since(start))
}
