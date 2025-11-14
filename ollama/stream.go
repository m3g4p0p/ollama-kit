package ollama

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"errors"
	"io"
	"iter"
)

type ResponseStream struct {
	body    io.ReadCloser
	decoder *jsontext.Decoder
	err     error
}

func NewResponseStream(in io.ReadCloser) *ResponseStream {
	return &ResponseStream{
		body:    in,
		decoder: jsontext.NewDecoder(in),
	}
}

func (r *ResponseStream) Err() error {
	if errors.Is(r.err, io.EOF) {
		return nil
	}

	return r.err
}

func (r *ResponseStream) Close() error {
	return r.body.Close()
}

func (r *ResponseStream) Recv() (ChatResponse, error) {
	var part ChatResponse

	if r.err != nil {
		return part, r.err
	}

	err := json.UnmarshalDecode(r.decoder, &part)
	if err != nil {
		r.err = err
	}

	return part, err
}

func (r *ResponseStream) Iter() iter.Seq[ChatResponse] {
	return func(yield func(ChatResponse) bool) {
		for {
			part, err := r.Recv()
			if err != nil || !yield(part) {
				return
			}
		}
	}
}

func (r *ResponseStream) Accumulate() iter.Seq2[ChatResponse, ChatMessage] {
	return func(yield func(ChatResponse, ChatMessage) bool) {
		var msg ChatMessage

		for {
			part, err := r.Recv()
			if err != nil {
				return
			}

			msg.Role = part.Message.Role
			msg.Content += part.Message.Content
			msg.ToolCalls = append(msg.ToolCalls, part.Message.ToolCalls...)

			if !yield(part, msg) {
				return
			}
		}
	}
}
