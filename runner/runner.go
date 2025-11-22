package runner

import (
	"context"

	"m3g4p0p/agents/ollama"
	"m3g4p0p/agents/runner/history"
	"m3g4p0p/agents/toolset"
	"m3g4p0p/agents/util"
)

type Runner struct {
	client    ollama.Client
	toolset   toolset.Toolset
	processor history.Processor
}

func NewRunner(client ollama.Client, options ...util.Option[Runner]) Runner {
	r := Runner{client: client}
	return util.ApplyOptions(r, options)
}

func (r Runner) Run(ctx context.Context, chat ollama.ChatRequest) *RunResult {
	chat.Tools = append(chat.Tools, r.toolset.Tools()...)

	if r.processor != nil {
		chat.Messages = r.processor(chat.Messages)
	}

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
	options ...util.Option[Runner],
) *RunResult {
	return NewRunner(client, options...).Run(ctx, chat)
}
