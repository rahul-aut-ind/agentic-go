package main

import (
	"context"
	"fmt"
	"time"

	"github.com/rahul-aut-ind/genkit-go/pkg/config"
	"github.com/rahul-aut-ind/genkit-go/pkg/logger"
	"github.com/rahul-aut-ind/genkit-go/services/imagegeneratorservice"
)

var (
	concept = "Developer coding in golang"
)

func main() {
	// initialize everything
	env := config.NewConfig()
	log := logger.New()

	service := imagegeneratorservice.NewClient(context.Background(), env, log)

	fmt.Println("Enter concept for which you want an image to be generated:")
	_, e := fmt.Scanf("%s", &concept)
	if e != nil {
		log.Warnf("Error reading concept: %w, taking default", e)
	}
	fmt.Println("ok, generating a image for you based on the concept : ", concept)

	// generate image
	start := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	resp, err := service.GenerateImage(ctx, concept)
	if err != nil {
		log.Errorf("Error generating image: %v", err)
		return
	}

	log.PrettifyJSON(resp)
	log.Infof("Image generated in %s", time.Since(start))

}
