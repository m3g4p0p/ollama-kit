package runner

import (
	"encoding/json"

	"m3g4p0p/agents/ollama"

	"github.com/google/jsonschema-go/jsonschema"
)

type handlerFunc func([]byte) string

type Toolset struct {
	tools    []ollama.Tool
	handlers map[string]handlerFunc
}

func NewToolset(options ...ToolOption) *Toolset {
	toolset := &Toolset{handlers: make(map[string]handlerFunc)}

	for _, opt := range options {
		opt(toolset)
	}

	return toolset
}

func AddTool[T any](t *Toolset, name, description string, handler func(T) string) error {
	schema, err := jsonschema.For[T](&jsonschema.ForOptions{})
	if err != nil {
		return err
	}

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

	return nil
}

type ToolOption func(t *Toolset)

func WithTool[T any](name, description string, handler func(T) string) ToolOption {
	return func(t *Toolset) {
		if err := AddTool(t, name, description, handler); err != nil {
			panic(err)
		}
	}
}
