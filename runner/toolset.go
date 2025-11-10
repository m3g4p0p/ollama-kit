package runner

import (
	"encoding/json"

	"m3g4p0p/agents/ollama"

	"github.com/google/jsonschema-go/jsonschema"
)

type Toolset struct {
	tools    []ollama.Tool
	handlers map[string]func([]byte) string
}

type ToolOption func(t *Toolset)

func WithTool[T any](name, description string, handler func(T) string) ToolOption {
	schema, err := jsonschema.For[T](&jsonschema.ForOptions{})
	if err != nil {
		panic(err)
	}

	return func(t *Toolset) {
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
			err := json.Unmarshal(raw, &args)
			if err != nil {
				return err.Error()
			}
			return handler(args)
		}
	}
}
