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
	return func(yield func(ollama.ChatResponse) bool) {
		for {
			for part := range r.RunResult.Stream() {
				if !yield(part) {
					return
				}
			}

			handoff, ok := r.Output().(HandoffParams)
			if !ok {
				return
			}

			*r = *handoff.agent.Run(r.Context(), r.Messages())
		}
	}
}
