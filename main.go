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
		pretty bool
		raw    bool
	}

	flag.BoolVar(&options.pretty, "pretty", false, "")
	flag.BoolVar(&options.raw, "raw", false, "")

	var chat ollama.ChatRequest

	flag.StringVar(&chat.Model, "model", "qwen3:1.7b", "")
	flag.BoolVar(&chat.Think, "think", true, "")
	flag.Parse()

	client := ollama.Client{BaseURL: "http://localhost:11434"}

	chat.Messages = append(chat.Messages, ollama.ChatMessage{
		Role:    "user",
		Content: "Say hello world!",
	})

	stream, err := client.Chat(context.Background(), chat)
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
