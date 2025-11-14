package main

import (
	"context"
	"flag"
	"log"

	"m3g4p0p/agents/console"
	"m3g4p0p/agents/ollama"
	"m3g4p0p/agents/runner"
	"m3g4p0p/agents/tools"
	"m3g4p0p/agents/toolset"

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

	for _, prompt := range fs.Args() {
		chat.Messages = append(chat.Messages, ollama.ChatMessage{
			Role:    "user",
			Content: prompt,
		})

		run := runner.Run(
			context.Background(),
			client,
			chat,
			toolset.WithTool(
				"get_weather",
				"Get the weather for the provided location",
				tools.GetWeather,
			),
		)

		console.WriteStream(
			run.Stream(),
			options.raw,
			options.pretty,
		)

		if err := run.Err(); err != nil {
			log.Fatal(err)
		}

		chat.Messages = run.Messages()
	}
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
