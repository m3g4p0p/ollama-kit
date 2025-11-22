package console

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"fmt"
	"io"
	"iter"
	"strings"

	"m3g4p0p/agents/ollama"
	"m3g4p0p/agents/util"

	"github.com/alecthomas/chroma/quick"
)

const defaultStyle = "fruity"

type Console struct {
	out    io.Writer
	style  string
	pretty bool
}

func NewConsole(out io.Writer, options ...util.Option[Console]) *Console {
	console := &Console{out: out, style: defaultStyle}
	util.ApplyOptionsTo(console, options)
	return console
}

func (c *Console) WriteJSON(v any) error {
	defer fmt.Fprintln(c.out)

	if c.pretty {
		return c.writePretty(v)
	}

	return c.writeRaw(v)
}

func (c *Console) WriteStream(stream iter.Seq2[ollama.ChatResponse, error], raw bool) error {
	for part, err := range stream {
		if err != nil {
			return err
		}

		if raw {
			c.WriteJSON(part)
		} else {
			c.writePart(part)
		}
	}

	return nil
}

func (c *Console) writeRaw(v any) error {
	return json.MarshalWrite(c.out, v)
}

func (c *Console) writePretty(v any) error {
	var s strings.Builder

	err := json.MarshalWrite(&s, v, jsontext.WithIndent("  "))
	if err != nil {
		return err
	}

	return quick.Highlight(c.out, s.String(), "json", "terminal256", c.style)
}

func (c *Console) writePart(part ollama.ChatResponse) {
	if part.Message.Content != "" {
		fmt.Fprintf(c.out, "\033[1m%s\033[0m", part.Message.Content)
	}
	if part.Message.Thinking != "" {
		fmt.Fprint(c.out, part.Message.Thinking)
	}
}
