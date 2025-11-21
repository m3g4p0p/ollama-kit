package agentlib

import (
	"iter"

	"m3g4p0p/agents/ollama"
	"m3g4p0p/agents/runner"
)

type AgentRunResult struct {
	*runner.RunResult
}

func (r *AgentRunResult) Stream() iter.Seq2[ollama.ChatResponse, error] {
	return func(yield func(ollama.ChatResponse, error) bool) {
		for {
			for part, err := range r.RunResult.Stream() {
				if !yield(part, err) || err != nil {
					return
				}
			}

			handoff, ok := r.Output().(Agent)
			if !ok {
				return
			}

			*r = *handoff.Run(r.Context(), r.Messages())
		}
	}
}
