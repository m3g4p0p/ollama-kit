package toolset

import (
	"context"
	"encoding/json/v2"
	"os/exec"
	"strings"

	"m3g4p0p/agents/ollama"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type MCPToolset struct {
	session *mcp.ClientSession
}

func NewMCPToolset(session *mcp.ClientSession) *MCPToolset {
	return &MCPToolset{session: session}
}

func (t *MCPToolset) Tools(ctx context.Context) ([]ollama.Tool, error) {
	res, err := t.session.ListTools(ctx, &mcp.ListToolsParams{})
	if err != nil {
		return nil, err
	}

	var tools []ollama.Tool

	for _, t := range res.Tools {
		tools = append(tools, ollama.Tool{
			Type: "function",
			Function: ollama.ToolFunction{
				Name:        t.Name,
				Description: t.Description,
				Parameters:  t.InputSchema,
			},
		})
	}

	return tools, nil
}

func (t *MCPToolset) Handle(ctx context.Context, call ollama.ToolCall) (ToolResult, error) {
	res, err := t.session.CallTool(ctx, &mcp.CallToolParams{
		Name:      call.Function.Name,
		Arguments: call.Function.Arguments,
	})
	if err != nil {
		return ToolResult{}, err
	}

	var s strings.Builder
	for _, c := range res.Content {
		err := json.MarshalWrite(&s, c)
		if err != nil {
			return ToolResult{}, err
		}
	}

	return ToolResult{
		Args:    res.StructuredContent,
		Content: s.String(),
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
