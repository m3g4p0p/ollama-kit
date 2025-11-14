package ollama

import "encoding/json"

type ChatConfig struct {
	Model   string         `json:"model"`
	Think   bool           `json:"think"`
	Stream  bool           `json:"stream"`
	Options map[string]any `json:"options,omitempty"`
}

type ChatMessage struct {
	Role      string     `json:"role"`
	Content   string     `json:"content"`
	ToolCalls []ToolCall `json:"tool_calls,omitempty"`
}

type ChatRequest struct {
	ChatConfig
	Messages []ChatMessage `json:"messages"`
	Tools    []Tool        `json:"tools,omitempty"`
}

type ChatResponse struct {
	Model     string `json:"model"`
	CreatedAt string `json:"created_at"`
	Message   struct {
		ChatMessage
		Thinking string `json:"thinking,omitempty"`
	} `json:"message"`
	Done bool `json:"done"`
}

type Tool struct {
	Type     string       `json:"type"`
	Function ToolFunction `json:"function"`
}

type ToolFunction struct {
	Name        string `json:"name"`
	Parameters  any    `json:"parameters"`
	Description string `json:"description,omitempty"`
}

type ToolCall struct {
	Function struct {
		Name        string          `json:"name"`
		Description string          `json:"description,omitempty"`
		Arguments   json.RawMessage `json:"arguments,omitempty"`
	} `json:"function"`
}
