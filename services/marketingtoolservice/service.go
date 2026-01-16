package marketingtoolservice

import (
	"context"
	"fmt"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/core"
	"github.com/firebase/genkit/go/genkit"
	"github.com/firebase/genkit/go/plugins/googlegenai"
	"github.com/rahul-aut-ind/genkit-go/domain/models"
	"github.com/rahul-aut-ind/genkit-go/pkg/config"
	"github.com/rahul-aut-ind/genkit-go/pkg/logger"
	"golang.org/x/sync/errgroup"
	"google.golang.org/genai"
)

const (
	ParallelToolFlowName = "parallelToolFlow"
	WeatherToolName      = "getWeather"

	WeatherToolDefinition = `Gets the current weather for a specific location`
)

type (
	Client struct {
		genkit       *genkit.Genkit
		log          *logger.Logger
		weatherTool  ai.Tool
		parallelFlow *core.Flow[*models.MarketingIdeaRequest, *models.MarketingIdeaResponse, struct{}]
	}
)

func NewClient(ctx context.Context, cfg *config.Config, log *logger.Logger) *Client {
	c := &Client{
		genkit: genkit.Init(ctx,
			genkit.WithPlugins(&googlegenai.GoogleAI{
				APIKey: cfg.APIKey,
			}),
			genkit.WithDefaultModel(cfg.GeneralModelName),
		),
		log: log,
	}
	c.weatherTool = genkit.DefineTool(c.genkit,
		WeatherToolName,
		WeatherToolDefinition,
		c.getWeatherTool,
	)
	c.parallelFlow = genkit.DefineFlow(c.genkit, ParallelToolFlowName, c.generateParallelFlow)
	return c
}

func (c *Client) GenerateMarketingIdea(ctx context.Context, productName, location string) (*models.MarketingIdeaResponse, error) {
	// Run the flow
	resp, err := c.parallelFlow.Run(ctx, &models.MarketingIdeaRequest{
		ProductName: productName,
		Location:    location,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to generate marketing idea: %v", err)
	}
	return resp, nil
}

func (c *Client) generateParallelFlow(ctx context.Context, req *models.MarketingIdeaRequest) (*models.MarketingIdeaResponse, error) {

	var name, tagline, weather string

	// Task 0: Get weather for location
	resp, err := genkit.Generate(ctx, c.genkit,
		ai.WithPrompt(models.GetWeatherPrompt, req.Location),
		ai.WithTools(c.weatherTool),
	)
	if err != nil {
		weather = "not available"
	} else {
		weather = resp.Text()
	}

	c.log.Infof("Weather in %s: %s", req.Location, weather)

	g, ctx := errgroup.WithContext(ctx)

	// Task 1: Generate a creative name
	g.Go(func() error {
		resp, err := genkit.Generate(ctx, c.genkit,
			ai.WithPrompt(models.CreativeNamePrompt, req.ProductName, req.Location, weather),
			ai.WithConfig(&genai.GenerateContentConfig{
				Temperature: genai.Ptr[float32](0.4),
				TopP:        genai.Ptr[float32](0.8),
				TopK:        genai.Ptr[float32](10.0),
			}),
		)
		if err != nil {
			return err
		}
		name = resp.Text()
		return nil
	})

	// Task 2: Generate a catchy tagline
	g.Go(func() error {
		resp, err := genkit.Generate(ctx, c.genkit,
			ai.WithPrompt(models.CreativeTaglinePrompt, req.ProductName, req.Location, weather),
			ai.WithConfig(&genai.GenerateContentConfig{
				Temperature: genai.Ptr[float32](0.5),
				TopP:        genai.Ptr[float32](0.4),
				TopK:        genai.Ptr[float32](20.0),
			}),
		)
		if err != nil {
			return err
		}
		tagline = resp.Text()
		return nil
	})

	if err := g.Wait(); err != nil {
		return nil, fmt.Errorf("failed to generate marketing copy: %v", err)
	}

	return &models.MarketingIdeaResponse{
		Name:    name,
		Tagline: tagline,
	}, nil
}

//nolint:revive // signature is important for genkit
func (c *Client) getWeatherTool(ctx *ai.ToolContext, input *models.WeatherToolInput) (string, error) {
	c.log.Infof("Tool is called for location: %s", input.Location)
	return fmt.Sprintf("The current weather in %s is 75°F and sunny.", input.Location), nil
}
