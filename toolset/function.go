package toolset

import (
	"context"
	"encoding/json"

	"m3g4p0p/agents/ollama"

	"github.com/google/jsonschema-go/jsonschema"
)

var NewToolset = NewFunctionToolset

type (
	handlerFunc          func(context.Context, []byte) (ToolResult, error)
	Handler[T any]       func(context.Context, T) (ToolResult, error)
	SimpleHandler[T any] func(T) (string, error)
)

type FunctionToolset struct {
	tools    []ollama.Tool
	handlers map[string]handlerFunc
}

func NewFunctionToolset(options ...FunctionToolsetOption) *FunctionToolset {
	toolset := &FunctionToolset{handlers: make(map[string]handlerFunc)}

	for _, opt := range options {
		opt(toolset)
	}

	return toolset
}

func (t *FunctionToolset) Tools() []ollama.Tool {
	return t.tools
}

func (t *FunctionToolset) Handle(ctx context.Context, call ollama.ToolCall) (ToolResult, error) {
	handler := t.handlers[call.Function.Name]
	return handler(ctx, call.Function.Arguments)
}

func AddTool[T any](t *FunctionToolset, name, description string, handler Handler[T]) error {
	err := addToolDef[T](t, name, description)
	if err != nil {
		return err
	}

	t.handlers[name] = func(ctx context.Context, raw []byte) (ToolResult, error) {
		var args T
		err := json.Unmarshal(raw, &args)
		if err != nil {
			return ToolResult{}, err
		}

		return handler(ctx, args)
	}

	return nil
}

func AddSimpleTool[T any](t *FunctionToolset, name, description string, handler SimpleHandler[T]) error {
	return AddTool(t, name, description, func(ctx context.Context, args T) (ToolResult, error) {
		content, err := handler(args)
		return ToolResult{Args: args, Content: content}, err
	})
}

func AddStructuredOutput[T any](t *FunctionToolset, name, description string) error {
	err := addToolDef[T](t, name, description)
	if err != nil {
		return err
	}

	t.handlers[name] = func(ctx context.Context, raw []byte) (ToolResult, error) {
		var args T
		err := json.Unmarshal(raw, &args)
		if err != nil {
			return ToolResult{}, err
		}

		result := ToolResult{Final: true, Args: args}
		return result, nil
	}

	return nil
}

func addToolDef[T any](t *FunctionToolset, name, description string) error {
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

	return nil
}
