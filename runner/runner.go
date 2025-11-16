package runner

import (
	"context"

	"m3g4p0p/agents/ollama"
	"m3g4p0p/agents/toolset"
)

type Runner struct {
	client  ollama.Client
	toolset toolset.Toolset
}

func NewRunner(client ollama.Client, toolset toolset.Toolset) Runner {
	return Runner{client: client, toolset: toolset}
}

func (r Runner) Run(ctx context.Context, chat ollama.ChatRequest) *RunResult {
	return &RunResult{
		ctx:     ctx,
		client:  r.client,
		toolset: r.toolset,
		chat:    chat,
	}
}

func Run(
	ctx context.Context,
	client ollama.Client,
	chat ollama.ChatRequest,
	toolset toolset.Toolset,
) *RunResult {
	return NewRunner(client, toolset).Run(ctx, chat)
}
