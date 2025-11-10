package runner

import (
	"encoding/json"

	"m3g4p0p/agents/ollama"

	"github.com/google/jsonschema-go/jsonschema"
)

type toolset struct {
	tools    []ollama.Tool
	handlers map[string]func([]byte) string
}

type toolOption func(t *toolset)

func WithTool[T any](name, description string, handler func(T) string) toolOption {
	schema, err := jsonschema.For[T](&jsonschema.ForOptions{})
	if err != nil {
		panic(err)
	}

	return func(t *toolset) {
		t.tools = append(t.tools, ollama.Tool{
			Type: "function",
			Function: ollama.ToolFunction{
				Name:        name,
				Description: description,
				Parameters:  schema,
			},
		})

		t.handlers[name] = func(raw []byte) string {
			var args T
			json.Unmarshal(raw, &args)
			return handler(args)
		}
	}
}
