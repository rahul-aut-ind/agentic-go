package main

import (
	"context"
	"fmt"
	"time"

	"github.com/rahul-aut-ind/genkit-go/pkg/config"
	"github.com/rahul-aut-ind/genkit-go/pkg/logger"
	"github.com/rahul-aut-ind/genkit-go/services/bloggeneratorservice"
)

var (
	topic = "The Future of AI"
)

func main() {
	// initialize everything
	env := config.NewConfig()
	log := logger.New()

	service := bloggeneratorservice.NewClient(context.Background(), env, log)

	fmt.Println("Type your topic:")
	_, e := fmt.Scanf("%s", &topic)
	if e != nil {
		log.Warnf("Error reading topic: %w, taking default", e)
	}

	fmt.Printf("\nok, generating a blog post for you based on the topic %s\n", topic)

	// generate blog and iteratively refine it
	start := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	resp, err := service.GenerateBlog(ctx, topic)
	if err != nil {
		log.Errorf("Error generating response: %v", err)
		return
	}

	log.PrettifyJSON(resp)
	log.Infof("Response generated in %s", time.Since(start))
}
