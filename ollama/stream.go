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
}

type Accumulation struct {
	Part    ChatResponse
	Message ChatMessage
}

func NewResponseStream(in io.ReadCloser) *ResponseStream {
	return &ResponseStream{
		body:    in,
		decoder: jsontext.NewDecoder(in),
	}
}

func (r *ResponseStream) Close() error {
	return r.body.Close()
}

func (r *ResponseStream) Recv() (ChatResponse, error) {
	var part ChatResponse
	err := json.UnmarshalDecode(r.decoder, &part)
	return part, err
}

func (r *ResponseStream) Iter() iter.Seq2[ChatResponse, error] {
	return func(yield func(ChatResponse, error) bool) {
		for {
			part, err := r.Recv()
			if errors.Is(err, io.EOF) {
				return
			}

			if !yield(part, err) || err != nil {
				return
			}
		}
	}
}

func (r *ResponseStream) Accumulate() iter.Seq2[Accumulation, error] {
	return func(yield func(Accumulation, error) bool) {
		var msg ChatMessage

		for part, err := range r.Iter() {
			if err != nil {
				yield(Accumulation{}, err)
				return
			}

			msg.Role = part.Message.Role
			msg.Content += part.Message.Content
			msg.ToolCalls = append(msg.ToolCalls, part.Message.ToolCalls...)

			if !yield(Accumulation{part, msg}, err) {
				return
			}
		}
	}
}
