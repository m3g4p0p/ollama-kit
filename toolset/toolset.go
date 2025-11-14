package toolset

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

func (t *Toolset) Tools() []ollama.Tool {
	return t.tools
}

func (t *Toolset) Handle(name string, rawArgs []byte) string {
	return t.handlers[name](rawArgs)
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
