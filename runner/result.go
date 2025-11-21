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
	output  any
}

func (r *RunResult) Output() any {
	return r.output
}

func (r *RunResult) Messages() []ollama.ChatMessage {
	return r.chat.Messages
}

func (r *RunResult) Stream() iter.Seq2[ollama.ChatResponse, error] {
	return func(yield func(ollama.ChatResponse, error) bool) {
		r.doRun(r.ctx, yield)
	}
}

func (r *RunResult) doRun(ctx context.Context, yield func(ollama.ChatResponse, error) bool) {
	for {
		stream, err := r.client.Chat(ctx, r.chat)
		if err != nil {
			yield(ollama.ChatResponse{}, err)
			return
		}
		defer stream.Close()

		for acc, err := range stream.Accumulate() {
			if !yield(acc.Part, err) || err != nil {
				return
			}

			if !acc.Part.Done {
				continue
			}

			r.chat.Messages = append(r.chat.Messages, acc.Message)

			if len(acc.Message.ToolCalls) == 0 {
				return
			}

			for _, call := range acc.Message.ToolCalls {
				var content string

				result, err := r.toolset.Handle(r.ctx, call)
				if err != nil {
					content = fmt.Sprintf("Error: %v", err)
				} else {
					content = result.Content
				}

				if result.Final && err == nil {
					r.output = result.Args
					return
				}

				r.chat.Messages = append(
					r.chat.Messages,
					ollama.ChatMessage{Role: "tool", Content: content},
				)
			}
		}

		stream.Close()
	}
}
