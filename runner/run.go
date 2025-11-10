package runner

import (
	"context"
	"iter"

	"m3g4p0p/agents/ollama"
)

type RunStream struct {
	client  ollama.Client
	chat    ollama.ChatRequest
	toolset Toolset
	err     error
}

func (r *RunStream) Err() error {
	return r.err
}

func (r *RunStream) Stream(ctx context.Context) iter.Seq[ollama.ChatResponse] {
	return func(yield func(ollama.ChatResponse) bool) {
		for {
			stream, err := r.client.Chat(ctx, r.chat)
			if err != nil {
				r.err = err
				return
			}
			defer stream.Close()

			var message ollama.ChatMessage
			var toolResults []string

			for part := range stream.Iter() {
				if !yield(part) {
					return
				}

				message.Role = part.Message.Role
				message.Content += part.Message.Content

				if part.Message.ToolCalls == nil {
					continue
				}

				for _, call := range part.Message.ToolCalls {
					handler := r.toolset.handlers[call.Function.Name]
					content := handler(call.Function.Arguments)
					toolResults = append(toolResults, content)
				}

				if !part.Done {
					continue
				}

				r.chat.Messages = append(r.chat.Messages, message)

				if len(toolResults) == 0 {
					return
				}

				for _, result := range toolResults {
					r.chat.Messages = append(
						r.chat.Messages,
						ollama.ChatMessage{Role: "tool", Content: result},
					)
				}
			}
		}
	}
}

func Run(
	client ollama.Client,
	chat ollama.ChatRequest,
	options ...ToolOption,
) *RunStream {
	var toolset Toolset

	for _, opt := range options {
		opt(&toolset)
	}

	chat.Tools = append(chat.Tools, toolset.tools...)

	return &RunStream{
		client:  client,
		chat:    chat,
		toolset: toolset,
	}
}
