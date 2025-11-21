package main

import (
	"context"
	"log"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type SayHiInput struct {
	Name string `json:"name" jsonschema:"the name of the person to greet"`
}

type SayHiOutput struct {
	Greeting string `json:"greeting" jsonschema:"the greeting to tell to the user"`
}

func SayHi(ctx context.Context, req *mcp.CallToolRequest, input SayHiInput) (
	*mcp.CallToolResult,
	SayHiOutput,
	error,
) {
	return nil, SayHiOutput{Greeting: "Hi " + input.Name}, nil
}

type GetWeatherInput struct {
	Location string `json:"location" jsonschema:"the location to get the weather for"`
}

type GetWeatherOutput struct {
	Condition string `json:"condition" jsonschema:"the weather condition"`
}

func GetWeather(ctx context.Context, req *mcp.CallToolRequest, input GetWeatherInput) (
	*mcp.CallToolResult,
	GetWeatherOutput,
	error,
) {
	res, err := req.Session.CreateMessage(ctx, &mcp.CreateMessageParams{
		Messages: []*mcp.SamplingMessage{{
			Role: "user",
			Content: &mcp.TextContent{
				Text: "give me a random weather condition",
			},
		}},
	})
	if err != nil {
		return nil, GetWeatherOutput{}, err
	}

	text := res.Content.(*mcp.TextContent).Text
	return nil, GetWeatherOutput{Condition: text}, nil
}

func server(args []string) {
	server := mcp.NewServer(&mcp.Implementation{Name: "greeter", Version: "v1.0.0"}, nil)
	mcp.AddTool(server, &mcp.Tool{Name: "greet", Description: "say hi"}, SayHi)
	mcp.AddTool(server, &mcp.Tool{Name: "get_weather", Description: "get the weather"}, GetWeather)

	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatal(err)
	}
}
