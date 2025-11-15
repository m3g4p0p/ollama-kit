package toolset

import (
	"encoding/json"

	"m3g4p0p/agents/ollama"

	"github.com/google/jsonschema-go/jsonschema"
)

type (
	handlerFunc    func([]byte) (ToolResult, error)
	Handler[T any] func(T) (string, error)
)

type ToolResult struct {
	Content string
	Final   bool
}

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

func (t Toolset) Handle(call ollama.ToolCall) (ToolResult, error) {
	handler := t.handlers[call.Function.Name]
	return handler(call.Function.Arguments)
}

func AddTool[T any](t *Toolset, name, description string, handler Handler[T]) error {
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

	t.handlers[name] = func(raw []byte) (ToolResult, error) {
		var args T
		var tr ToolResult

		err := json.Unmarshal(raw, &args)
		if err != nil {
			return tr, err
		}

		tr.Content, err = handler(args)
		return tr, err
	}

	return nil
}
