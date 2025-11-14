package main

import (
	"context"
	"flag"
	"log"

	"m3g4p0p/agents/agentlib"
	"m3g4p0p/agents/console"
	"m3g4p0p/agents/ollama"
	"m3g4p0p/agents/tools"

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

	agent := agentlib.NewAgent(
		ollama.Client{BaseURL: "http://localhost:11434"},
		agentlib.WithChatRequest(chat),
		agentlib.WithTool(
			"get_weather",
			"Get the weather for the provided location",
			func(params tools.GetWeatherParams) (string, error) {
				return "sunny", nil
			},
		),
	)

	run := agent.Run(context.Background(), fs.Arg(0))

	console.WriteStream(
		run.Stream(),
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
