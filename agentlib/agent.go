package agentlib

import (
	"context"

	"m3g4p0p/agents/agentlib/history"
	"m3g4p0p/agents/ollama"
	"m3g4p0p/agents/runner"
	"m3g4p0p/agents/toolset"
)

var defaultProcessor = history.Combine(
	history.DeleteRoles("system", "tool"),
	history.DeleteToolCalls(),
)

type Agent struct {
	client    ollama.Client
	chat      ollama.ChatRequest
	toolset   toolset.Toolset
	processer history.Processor
}

func NewAgent(client ollama.Client, options ...Option) Agent {
	agent := Agent{
		client:    client,
		toolset:   toolset.NewToolset(),
		processer: defaultProcessor,
	}

	for _, opt := range options {
		opt(&agent)
	}

	return agent
}

func (a Agent) Run(ctx context.Context, prompt string, history []ollama.ChatMessage) *runner.RunResult {
	msg := ollama.ChatMessage{
		Role:    "user",
		Content: prompt,
	}

	history = append(a.processer(history), msg)
	a.chat.Messages = append(a.chat.Messages, history...)

	return runner.NewRunner(a.client, a.toolset).Run(ctx, a.chat)
}
