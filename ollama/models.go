package ollama

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Model    string         `json:"model"`
	Messages []ChatMessage  `json:"messages"`
	Think    *bool          `json:"think,omitempty"`
	Stream   *bool          `json:"stream,omitempty"`
	Options  map[string]any `json:"options,omitempty"`
}

type ChatResponse struct {
	Model     string `json:"model"`
	CreatedAt string `json:"created_at"`
	Message   struct {
		Role     string `json:"role"`
		Content  string `json:"content"`
		Thinking string `json:"thinking,omitempty"`
	} `json:"message"`
	Done bool `json:"done"`
}
