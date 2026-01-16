## Common hurdle with Genkit Go: the way tool input parameters are defined.

When we pass a raw string as the input argument to a tool function in Go, Genkit's reflection-based schema generator sometimes fails to provide a named property (like location) in the JSON schema it sends to Gemini. Without a named property, the model thinks the tool takes no arguments.

**The Solution: Use a Struct for Input**

In Genkit Go, we should almost always use a struct to define tool inputs. This allows Genkit to generate a proper JSON schema that the model can understand.
1. Define an Input Struct

Add a struct to our models or within our service to represent the tool's input:

```sh
type WeatherInput struct {
    Location string `json:"location" jsonschema_description:"The city and state, e.g. San Francisco, CA"`
}
```

2. Update the Tool Function Signature

Change the getWeatherTool function to accept this struct instead of a raw string:
Go
```sh
func (c *Client) getWeatherTool(ctx *ai.ToolContext, input *WeatherInput) (string, error) {
    c.log.Infof("Tool is called for location: %s", input.Location)
    return fmt.Sprintf("The current weather in %s is 75°F and sunny.", input.Location), nil
}
```

3. Improve the Tool Definition

When defining the tool, ensure the description is clear. Genkit uses the description to tell the model when and how to use it.
Go

```sh
c.weatherTool = genkit.DefineTool(c.genkit,
    WeatherToolName,
    "Gets the current weather for a specific location", // Clear description
    c.getWeatherTool,
)
```

Why this fixes it?

When we use a struct, Genkit generates a JSON Schema that looks like this behind the scenes:
JSON
```sh
{
  "name": "getWeather",
  "parameters": {
    "type": "object",
    "properties": {
      "location": { "type": "string" }
    },
    "required": ["location"]
  }
}
```

When we use just a **string** in the function signature, the schema generated is often just `{"type": "string"}`, which Gemini 2.5 sometimes interprets as a tool that doesn't accept structured arguments, leading to the "cannot fulfill this request" error we saw.

A Note on Tool Logic:

In our generateParallelFlow, we are calling resp.Text(). Keep in mind that when a tool is called, Genkit usually handles the "call-response-resubmit" loop automatically if we are using high-level interfaces, but since we are calling genkit.Generate manually:

    Gemini will return a Tool Call request.

    Genkit will execute our getWeatherTool.

    Genkit will then automatically send the tool output back to Gemini to get the final text response.

    resp.Text() will then contain the final natural language answer (e.g., "The weather in Berlin is 75°F...").

----

In Genkit Go, **"Automatic Tool Calling"** refers to the framework's ability to automatically execute a tool's logic and send the result back to the model without us having to write a manual loop.

This behavior is actually enabled by default whenever we provide tools to a genkit.Generate call. However, we can strictly control how the model chooses to use those tools using the ToolChoice configuration.
1. Enabling Forced Tool Calling

If we want to ensure the model always attempts to call a tool (rather than just answering with text), we use the ai.WithToolChoice option.

```sh
import "github.com/firebase/genkit/go/ai"

// ... inside our flow or function
resp, err := genkit.Generate(ctx, g,
    ai.WithModel("googleai/gemini-2.5-flash"),
    ai.WithPrompt(input.Question),
    ai.WithTools(c.weatherTool),
    // This ensures the model MUST use one of the provided tools
    ai.WithToolChoice(ai.ToolChoiceRequired), 
)
```
2. Available Tool Choice Options

Genkit provides three main modes for tool selection:
-Constant Behavior
- ai.ToolChoiceAuto	: Default. The model decides whether to call a tool or respond with text.
- ai.ToolChoiceRequired	: The model is forced to call at least one tool before returning a final response.
- ai.ToolChoiceNone : 	The model is prohibited from calling any tools (even if provided).

3. Understanding the "Automatic" Loop

It is important to distinguish between Selection Mode (above) and the Execution Loop.

By default, Genkit handles the following loop for us:

    Model requests a tool call.

    Genkit intercepts the request and runs our Go function.

    Genkit sends the tool's output back to the Model.

    Model generates the final natural language response.

    [!TIP] If we ever want to disable the automatic execution (e.g., if we need to ask a user for permission before a tool runs), we can add ai.WithReturnToolRequests(true). This will stop the loop and return the request to we to handle manually.

Summary of our updated Task

To make our weather tool robust and "forced," our generate call should look like this:

```sh
resp, err := genkit.Generate(ctx, g,
    ai.WithModel("googleai/gemini-2.5-flash"),
    ai.WithPrompt(input.Question),
    ai.WithTools(c.weatherTool),
    ai.WithToolChoice(ai.ToolChoiceRequired), // Force the interaction
)
```