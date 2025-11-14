package agentlib

import (
	"context"

	"m3g4p0p/agents/ollama"
	"m3g4p0p/agents/runner"
	"m3g4p0p/agents/toolset"
)

type Agent struct {
	client  ollama.Client
	chat    ollama.ChatRequest
	toolset toolset.Toolset
}

func NewAgent(client ollama.Client, options ...Option) Agent {
	agent := Agent{client: client, toolset: toolset.NewToolset()}

	for _, opt := range options {
		opt(&agent)
	}

	return agent
}

func (a Agent) Run(ctx context.Context, prompt string) *runner.RunStream {
	msg := ollama.ChatMessage{
		Role:    "user",
		Content: prompt,
	}

	a.chat.Messages = append(a.chat.Messages, msg)
	return runner.NewRunner(a.client, a.toolset).Run(ctx, a.chat)
}
