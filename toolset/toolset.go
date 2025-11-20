package toolset

import (
	"context"

	"m3g4p0p/agents/ollama"
)

type ToolResult struct {
	Args    any
	Content string
	Final   bool
}

type Toolset interface {
	Tools() []ollama.Tool
	Handle(ctx context.Context, call ollama.ToolCall) (ToolResult, error)
}
