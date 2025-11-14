package agent

import (
	"m3g4p0p/agents/ollama"
	"m3g4p0p/agents/toolset"
)

type Agent struct {
	client   ollama.Client
	chatBase ollama.ChatRequest
	toolset  toolset.Toolset
}

func NewAgent(client ollama.Client, options ...Option) Agent {
	agent := Agent{client: client}

	for _, opt := range options {
		opt(&agent)
	}

	return agent
}
