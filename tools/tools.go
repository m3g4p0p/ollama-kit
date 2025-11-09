package tools

import (
	"m3g4p0p/agents/ollama"

	"github.com/google/jsonschema-go/jsonschema"
)

type GetWeatherParams struct {
	Location string `json:"location"`
}

func AddTool[T any](chat *ollama.ChatRequest, name, description string) error {
	schema, err := jsonschema.For[T](&jsonschema.ForOptions{})
	if err != nil {
		return err
	}

	chat.Tools = append(chat.Tools, ollama.Tool{
		Type: "function",
		Function: ollama.ToolFunction{
			Name:        name,
			Description: description,
			Parameters:  schema,
		},
	})

	return nil
}
