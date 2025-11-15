package runner

import (
	"context"

	"m3g4p0p/agents/ollama"
	"m3g4p0p/agents/toolset"
)

type StructuredRunResult[T any] struct {
	*RunResult
}

func (r *StructuredRunResult[T]) Output() T {
	return r.RunResult.output.(T)
}

func RunStructured[T any](
	ctx context.Context,
	client ollama.Client,
	chat ollama.ChatRequest,
	toolset toolset.Toolset,
) *StructuredRunResult[T] {
	r := NewRunner(client, toolset).Run(ctx, chat)
	return &StructuredRunResult[T]{r}
}
