package toolset

import (
	"encoding/json"

	"m3g4p0p/agents/ollama"

	"github.com/google/jsonschema-go/jsonschema"
)

type handlerFunc func([]byte) (string, error)

type Toolset struct {
	tools    []ollama.Tool
	handlers map[string]handlerFunc
}

func NewToolset(options ...ToolOption) Toolset {
	toolset := Toolset{handlers: make(map[string]handlerFunc)}

	for _, opt := range options {
		opt(&toolset)
	}

	return toolset
}

func (t Toolset) Tools() []ollama.Tool {
	return t.tools
}

func (t Toolset) Handle(call ollama.ToolCall) (string, error) {
	handler := t.handlers[call.Function.Name]
	return handler(call.Function.Arguments)
}

func AddTool[T any](t *Toolset, name, description string, handler func(T) (string, error)) error {
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

	t.handlers[name] = func(raw []byte) (string, error) {
		var args T
		err := json.Unmarshal(raw, &args)
		if err != nil {
			return "", err
		}
		return handler(args)
	}

	return nil
}
