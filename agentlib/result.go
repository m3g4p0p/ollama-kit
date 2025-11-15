package agentlib

import (
	"iter"

	"m3g4p0p/agents/ollama"
	"m3g4p0p/agents/runner"
)

type AgentRunResult struct {
	*runner.RunResult
}

func (r *AgentRunResult) Stream() iter.Seq[ollama.ChatResponse] {
	return r.RunResult.Stream()
}
