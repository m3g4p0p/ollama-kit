package runner

import (
	"context"
	"fmt"
	"iter"

	"m3g4p0p/agents/ollama"
	"m3g4p0p/agents/toolset"
)

type RunResult struct {
	ctx     context.Context
	client  ollama.Client
	toolset toolset.Toolset
	chat    ollama.ChatRequest
	err     error
}

func (r *RunResult) Err() error {
	return r.err
}

func (r *RunResult) Messages() []ollama.ChatMessage {
	return r.chat.Messages
}

func (r *RunResult) Stream() iter.Seq[ollama.ChatResponse] {
	return func(yield func(ollama.ChatResponse) bool) {
		r.doRun(r.ctx, yield)
	}
}

func (r *RunResult) doRun(ctx context.Context, yield func(ollama.ChatResponse) bool) {
	for {
		stream, err := r.client.Chat(ctx, r.chat)
		if err != nil {
			r.err = err
			return
		}
		defer stream.Close()

		for part, message := range stream.Accumulate() {
			if !yield(part) {
				return
			}

			if !part.Done {
				continue
			}

			r.chat.Messages = append(r.chat.Messages, message)

			if len(message.ToolCalls) == 0 {
				return
			}

			for _, call := range message.ToolCalls {
				result, err := r.toolset.Handle(call)
				if err != nil {
					result = fmt.Sprintf("Error: %v", err)
				}

				r.chat.Messages = append(
					r.chat.Messages,
					ollama.ChatMessage{Role: "tool", Content: result},
				)
			}
		}

		if err := stream.Err(); err != nil {
			r.err = err
			return
		}

		stream.Close()
	}
}
