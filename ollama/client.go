package ollama

import (
	"context"
	"encoding/json/v2"
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	BaseURL string
}

func (c Client) Chat(ctx context.Context, chat ChatRequest) (*ResponseStream, error) {
	body, err := c.doRequest(ctx, "/api/chat", chat)
	if err != nil {
		return nil, err
	}

	return NewResponseStream(body), nil
}

func (c Client) doRequest(
	ctx context.Context,
	path string,
	data any,
) (io.ReadCloser, error) {
	r, w := io.Pipe()

	go func() {
		w.CloseWithError(json.MarshalWrite(w, data))
	}()

	req, err := http.NewRequestWithContext(ctx, "POST", c.BaseURL+path, r)
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

	return resp.Body, nil
}
