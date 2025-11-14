package runner

import (
	"context"
	"iter"

	"m3g4p0p/agents/ollama"
)

type Runner struct {
	client  ollama.Client
	toolset *Toolset
}

func NewRunner(client ollama.Client, toolset *Toolset) Runner {
	return Runner{client: client, toolset: toolset}
}

func (r Runner) Run(ctx context.Context, chat ollama.ChatRequest) *RunStream {
	chat.Tools = append(chat.Tools, r.toolset.tools...)

	return &RunStream{
		ctx:     ctx,
		client:  r.client,
		toolset: r.toolset,
		chat:    chat,
	}
}

type RunStream struct {
	ctx     context.Context
	client  ollama.Client
	chat    ollama.ChatRequest
	toolset *Toolset
	err     error
}

func (r *RunStream) Err() error {
	return r.err
}

func (r *RunStream) Stream() iter.Seq[ollama.ChatResponse] {
	return func(yield func(ollama.ChatResponse) bool) {
		r.doRun(r.ctx, yield)
	}
}

func (r *RunStream) doRun(ctx context.Context, yield func(ollama.ChatResponse) bool) {
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

func Run(
	ctx context.Context,
	client ollama.Client,
	chat ollama.ChatRequest,
	options ...ToolOption,
) *RunStream {
	return NewRunner(client, NewToolset(options...)).Run(ctx, chat)
}
