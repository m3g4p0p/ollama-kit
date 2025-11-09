package ollama

import (
	"bytes"
	"encoding/json/jsontext"
	"encoding/json/v2"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	BaseURL string
	Model   string
}

func (c Client) Chat(messages []ChatMessage) error {
	p, err := json.Marshal(ChatRequest{
		Model:    c.Model,
		Messages: messages,
	})
	if err != nil {
		return err
	}

	fmt.Println(c.BaseURL + "/api/chat")
	req, err := http.NewRequest("POST", c.BaseURL+"/api/chat", bytes.NewReader(p))
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		return fmt.Errorf("%s", resp.Status)
	}

	dec := jsontext.NewDecoder(resp.Body)

	for {
		part := ChatResponse{}
		err := json.UnmarshalDecode(dec, &part)
		if errors.Is(err, io.EOF) {
			return nil
		}

		if err != nil {
			return err
		}

		fmt.Printf("%+v\n", part)
	}
}
