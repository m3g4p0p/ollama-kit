package console

import (
	"encoding/json/v2"
	"fmt"
	"iter"
	"os"

	"m3g4p0p/agents/ollama"
	"m3g4p0p/agents/util"
)

func WriteStream(stream iter.Seq[ollama.ChatResponse], raw, pretty bool) {
	for part := range stream {
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
}
