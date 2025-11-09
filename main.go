package main

import (
	"context"
	"encoding/json/v2"
	"flag"
	"fmt"
	"log"
	"os"

	"m3g4p0p/agents/ollama"
	"m3g4p0p/agents/util"

	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	var options struct {
		model  string
		pretty bool
		raw    bool
	}

	flag.StringVar(&options.model, "model", "qwen3:1.7b", "")
	flag.BoolVar(&options.pretty, "pretty", false, "")
	flag.BoolVar(&options.raw, "raw", false, "")
	flag.Parse()

	client := ollama.Client{
		BaseURL: "http://localhost:11434",
		Model:   options.model,
	}

	stream, err := client.Chat(context.Background(), []ollama.ChatMessage{
		{Role: "user", Content: "Say hello world!"},
	})
	if err != nil {
		log.Fatal(err)
	}

	for part := range stream.Iter() {
		if options.raw {
			if options.pretty {
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
