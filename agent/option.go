package agent

import (
	"m3g4p0p/agents/ollama"
	"m3g4p0p/agents/toolset"
)

type Option func(a *Agent)

func WithInstructions(instructions string) Option {
	return func(a *Agent) {
		a.chatBase.Messages = append(a.chatBase.Messages, ollama.ChatMessage{
			Role:    "system",
			Content: instructions,
		})
	}
}

func WithModel(model string) Option {
	return func(a *Agent) {
		a.chatBase.Model = model
	}
}

func WithStream(stream bool) Option {
	return func(a *Agent) {
		a.chatBase.Stream = stream
	}
}

func WithThink(think bool) Option {
	return func(a *Agent) {
		a.chatBase.Think = think
	}
}

func WithToolset(toolst toolset.Toolset) Option {
	return func(a *Agent) {
		a.toolset = toolst
	}
}
