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

func prompt(args []string) {
	var options struct {
		pretty bool
		raw    bool
	}

	fs := flag.NewFlagSet("prompt", flag.ExitOnError)
	fs.BoolVar(&options.pretty, "pretty", false, "")
	fs.BoolVar(&options.raw, "raw", false, "")

	var chat ollama.ChatRequest

	fs.StringVar(&chat.Model, "model", "qwen3:1.7b", "")
	fs.BoolVar(&chat.Stream, "stream", true, "")
	fs.BoolVar(&chat.Think, "think", false, "")
	fs.Parse(args)

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

var cmds = map[string]func([]string){
	"prompt": prompt,
}

func main() {
	flag.Parse()

	if flag.NArg() == 0 {
		log.Fatalln("missing command")
	}

	cmds[flag.Arg(0)](flag.Args()[1:])
}
