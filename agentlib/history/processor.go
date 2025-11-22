package history

import (
	"slices"

	"m3g4p0p/agents/ollama"
)

type Processor func(history []ollama.ChatMessage) []ollama.ChatMessage

func Combine(processors ...Processor) Processor {
	return func(history []ollama.ChatMessage) []ollama.ChatMessage {
		for _, p := range processors {
			history = p(history)
		}

		return history
	}
}

func Deleter(del func(msg ollama.ChatMessage) bool) Processor {
	return func(history []ollama.ChatMessage) []ollama.ChatMessage {
		return slices.DeleteFunc(history, del)
	}
}

func DeleteRole(role string) Processor {
	return Deleter(func(msg ollama.ChatMessage) bool {
		return msg.Role == role
	})
}

func DeleteRoles(roles ...string) Processor {
	return Deleter(func(msg ollama.ChatMessage) bool {
		return slices.Contains(roles, msg.Role)
	})
}

func DeleteToolCalls() Processor {
	return Deleter(func(msg ollama.ChatMessage) bool {
		return len(msg.ToolCalls) > 0
	})
}
