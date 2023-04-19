package main

import (
	"context"
	"flag"
	"github.com/patrickjmcd/openai-cli/pkg"
	openai "github.com/sashabaranov/go-openai"
	"log"
	"os"
)

var apiKey = ""
var chatGptVersion = openai.GPT3Dot5Turbo

func init() {
	if v, ok := os.LookupEnv("OPENAI_API_KEY"); ok {
		apiKey = v
	} else {
		log.Panic("OPENAI_API_KEY is not set")
	}

	versionFlag := flag.String("version", "3.5-turbo", "GPT3 version to use (3.5-turbo, 4)")
	flag.Parse()

	switch *versionFlag {
	case "3.5-turbo":
		chatGptVersion = openai.GPT3Dot5Turbo
	case "4":
		chatGptVersion = openai.GPT4
	default:
		log.Fatalf("Invalid version: %s", *versionFlag)
	}

}

func getJoinedPrompt() string {
	joinedPrompt := ""
	for i := 1; i < flag.NArg(); i++ {
		joinedPrompt += flag.Arg(i) + " "
	}
	return joinedPrompt
}

func main() {
	ctx := context.Background()
	client := openai.NewClient(apiKey)

	switch a := flag.Arg(0); a {
	case "ask", "chat", "chatgpt":
		pkg.ChatGPT(client, chatGptVersion, getJoinedPrompt())
	case "dalle", "img":
		pkg.DallE(ctx, client, getJoinedPrompt(), "")
	default:
		log.Fatalf("Invalid command: %s", a)
	}
}
