package console

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"fmt"
	"io"
	"iter"
	"os"
	"strings"

	"m3g4p0p/agents/ollama"
	"m3g4p0p/agents/util"

	"github.com/alecthomas/chroma/quick"
)

const defaultStyle = "fruity"

type Console struct {
	out   io.Writer
	style string
}

func NewConsole(out io.Writer, options ...util.Option[Console]) *Console {
	console := &Console{out: out, style: defaultStyle}
	util.ApplyOptionsTo(console, options)
	return console
}

func (c *Console) WriteJSON(v any) error {
	var s strings.Builder

	err := json.MarshalWrite(&s, v, jsontext.WithIndent("  "))
	if err != nil {
		return err
	}

	return quick.Highlight(c.out, s.String(), "json", "terminal256", c.style)
}

func (c *Console) WriteStream(stream iter.Seq2[ollama.ChatResponse, error], raw, pretty bool) error {
	for part, err := range stream {
		if err != nil {
			return err
		}

		if raw {
			if pretty {
				util.PrettyDumpJSON(os.Stdout, part)
			} else {
				json.MarshalWrite(os.Stdout, part)
			}

			fmt.Fprintln(os.Stdout)
		} else {
			if part.Message.Content != "" {
				fmt.Printf("\033[1m%s\033[0m", part.Message.Content)
			}
			if part.Message.Thinking != "" {
				fmt.Print(part.Message.Thinking)
			}
		}
	}

	return nil
}
