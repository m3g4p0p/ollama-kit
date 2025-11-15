package agentlib

import (
	"m3g4p0p/agents/agentlib/history"
	"m3g4p0p/agents/ollama"
	"m3g4p0p/agents/toolset"
)

type Option func(a *Agent)

func WithChatConfig(config ollama.ChatConfig) Option {
	return func(a *Agent) {
		a.chat.ChatConfig = config
	}
}

func WithInstructions(instructions string) Option {
	return func(a *Agent) {
		msg := ollama.ChatMessage{
			Role:    "system",
			Content: instructions,
		}

		a.chat.Messages = append(a.chat.Messages, msg)
	}
}

func WithHistoryProcessr(processor history.Processor) Option {
	return func(a *Agent) {
		a.processer = processor
	}
}

func WithModel(model string) Option {
	return func(a *Agent) {
		a.chat.Model = model
	}
}

func WithStream(stream bool) Option {
	return func(a *Agent) {
		a.chat.Stream = stream
	}
}

func WithThink(think bool) Option {
	return func(a *Agent) {
		a.chat.Think = think
	}
}

func WithTool[T any](name, description string, handler toolset.Handler[T]) Option {
	return func(a *Agent) {
		toolset.AddTool(&a.toolset, name, description, handler)
	}
}

func WithToolset(toolst toolset.Toolset) Option {
	return func(a *Agent) {
		a.toolset = toolst
	}
}
