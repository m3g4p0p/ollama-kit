package ollama

import (
	"bytes"
	"context"
	"encoding/json/v2"
	"fmt"
	"net/http"
)

type Client struct {
	BaseURL string
}

func (c Client) Chat(ctx context.Context, chat ChatRequest) (*ResponseStream, error) {
	p, err := json.Marshal(chat)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(
		ctx,
		"POST",
		c.BaseURL+"/api/chat",
		bytes.NewReader(p),
	)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		return nil, fmt.Errorf("%s", resp.Status)
	}

	return NewResponseStream(resp.Body), nil
}
