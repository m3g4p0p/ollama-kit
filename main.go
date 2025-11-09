package main

import (
	"context"
	"flag"
	"log"

	"m3g4p0p/agents/console"
	"m3g4p0p/agents/ollama"

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
	flag.BoolVar(&chat.Stream, "stream", true, "")
	flag.BoolVar(&chat.Think, "think", false, "")
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
	defer stream.Close()

	console.WriteStream(
		stream,
		options.raw,
		options.pretty,
	)
}
