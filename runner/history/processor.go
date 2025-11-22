package history

import (
	"slices"

	"m3g4p0p/agents/ollama"
)

type Processor func(history []ollama.ChatMessage) []ollama.ChatMessage

func AddInstructions(instrutions string) Processor {
	return func(history []ollama.ChatMessage) []ollama.ChatMessage {
		system := []ollama.ChatMessage{{Role: "system", Content: instrutions}}
		return append(system, history...)
	}
}

func Combine(processors ...Processor) Processor {
	processors = slices.DeleteFunc(processors, func(p Processor) bool {
		return p == nil
	})

	if len(processors) == 1 {
		return processors[0]
	}

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
