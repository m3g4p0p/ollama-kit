package main

import (
	"fmt"
	"log"

	"m3g4p0p/agents/ollama"

	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	client := ollama.Client{BaseURL: "http://localhost:11434", Model: "qwen3:1.7b"}
	stream, err := client.Chat([]ollama.ChatMessage{{Role: "user", Content: "Say hello world!"}})
	if err != nil {
		log.Fatal(err)
	}

	for part := range stream.Iter() {
		fmt.Print(part.Message.Content)
	}
}
