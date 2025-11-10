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

			for part := range stream.Iter() {
				if !yield(part) {
					return
				}

				message.Role = part.Message.Role
				message.Content += part.Message.Content
				message.ToolCalls = append(message.ToolCalls, part.Message.ToolCalls...)

				if !part.Done {
					continue
				}

				r.chat.Messages = append(r.chat.Messages, message)

				if len(message.ToolCalls) == 0 {
					return
				}

				for _, call := range message.ToolCalls {
					handler := r.toolset.handlers[call.Function.Name]
					result := handler(call.Function.Arguments)

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
		}
	}
}

func Run(
	client ollama.Client,
	chat ollama.ChatRequest,
	options ...ToolOption,
) *RunStream {
	toolset := Toolset{handlers: make(map[string]handlerFunc)}

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
