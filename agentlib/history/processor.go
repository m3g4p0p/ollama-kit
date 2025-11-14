package history

import (
	"slices"

	"m3g4p0p/agents/ollama"
)

type Processor func(history []ollama.ChatMessage) []ollama.ChatMessage

func RemoveRole(role string) Processor {
	return func(history []ollama.ChatMessage) []ollama.ChatMessage {
		return slices.DeleteFunc(history, func(msg ollama.ChatMessage) bool {
			return msg.Role == role
		})
	}
}

func RemoveRoles(roles ...string) Processor {
	return func(history []ollama.ChatMessage) []ollama.ChatMessage {
		return slices.DeleteFunc(history, func(msg ollama.ChatMessage) bool {
			return slices.Contains(roles, msg.Role)
		})
	}
}

func RemoveToolCalls() Processor {
	return func(history []ollama.ChatMessage) []ollama.ChatMessage {
		return slices.DeleteFunc(history, func(msg ollama.ChatMessage) bool {
			return len(msg.ToolCalls) > 0
		})
	}
}

func Combine(processors ...Processor) Processor {
	return func(history []ollama.ChatMessage) []ollama.ChatMessage {
		for _, p := range processors {
			history = p(history)
		}

		return history
	}
}
