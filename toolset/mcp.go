package toolset

import (
	"context"
	"os/exec"
	"strings"

	"m3g4p0p/agents/ollama"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type MCPToolset struct {
	tools   []ollama.Tool
	session *mcp.ClientSession
}

func NewMCPToolset(session *mcp.ClientSession) *MCPToolset {
	return &MCPToolset{session: session}
}

func (t *MCPToolset) Connect(ctx context.Context) error {
	for tool, err := range t.session.Tools(ctx, &mcp.ListToolsParams{}) {
		if err != nil {
			return err
		}

		t.tools = append(t.tools, ollama.Tool{
			Type: "function",
			Function: ollama.ToolFunction{
				Name:        tool.Name,
				Description: tool.Description,
				Parameters:  tool.InputSchema,
			},
		})
	}

	return nil
}

func (t *MCPToolset) Tools() []ollama.Tool {
	return t.tools
}

func (t *MCPToolset) Handle(ctx context.Context, call ollama.ToolCall) (ToolResult, error) {
	res, err := t.session.CallTool(ctx, &mcp.CallToolParams{
		Name:      call.Function.Name,
		Arguments: call.Function.Arguments,
	})
	if err != nil {
		return ToolResult{}, err
	}

	var parts []string
	for _, c := range res.Content {
		if tc, ok := c.(*mcp.TextContent); ok {
			parts = append(parts, tc.Text)
		}
	}

	return ToolResult{
		Args:    res.StructuredContent,
		Content: strings.Join(parts, "\n\n"),
	}, nil
}

func CreateCommandSession(ctx context.Context, cmd *exec.Cmd) (*mcp.ClientSession, error) {
	client := mcp.NewClient(&mcp.Implementation{
		Name:    "mcp-client",
		Version: "v1.0.0",
	}, nil)

	transport := &mcp.CommandTransport{Command: cmd}
	return client.Connect(ctx, transport, nil)
}
