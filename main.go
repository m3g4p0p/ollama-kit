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

	var config ollama.ChatConfig

	fs.StringVar(&config.Model, "model", "qwen3:1.7b", "")
	fs.BoolVar(&config.Stream, "stream", true, "")
	fs.BoolVar(&config.Think, "think", false, "")
	fs.Parse(args)

	agent := agentlib.NewAgent(
		ollama.Client{BaseURL: "http://localhost:11434"},
		agentlib.WithChatConfig(config),
		agentlib.WithSimpleTool(
			"get_weather",
			"Get the weather for the provided location",
			tools.GetWeather,
		),
	)

	var history []ollama.ChatMessage

	for _, arg := range fs.Args() {
		run := agent.Run(context.Background(), arg, history)

		console.WriteStream(
			run.Stream(),
			options.raw,
			options.pretty,
		)

		if err := run.Err(); err != nil {
			log.Fatal(err)
		}

		history = run.Messages()
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
