package main

import (
	"context"
	"flag"
	"log"
	"os"
	"strings"

	"github.com/sashabaranov/go-openai"
)

var (
	dependenciesCSV = flag.String("d", "", "dependency interfaces (comma separated)")
	specification   = flag.String("s", "", "specification interface")
	packageName     = flag.String("p", "", "package name")
	inputPath       = flag.String("i", "", "input path")
	outputPath      = flag.String("o", "", "output path")
)

func main() {
	flag.Parse()
	if *inputPath == "" || *outputPath == "" || *specification == "" || *packageName == "" {
		flag.Usage()
		return
	}
	key := os.Getenv("OPENAI_API_KEY_GENERATE_GPT")
	if key == "" {
		log.Fatal("OPENAI_API_KEY is not set")
	}
	client := openai.NewClient(key)
	dependencies := strings.Split(*dependenciesCSV, ",")
	if err := Generate(client, *inputPath, *outputPath, *packageName, *specification, dependencies); err != nil {
		log.Fatal(err)
	}
}

func Generate(client *openai.Client, in, out, pkg, spec string, deps []string) error {
	input, err := os.ReadFile(in)
	if err != nil {
		return err
	}
	resp, err := client.CreateChatCompletion(context.Background(), openai.ChatCompletionRequest{
		Model: openai.GPT3Dot5Turbo,
		Messages: []openai.ChatCompletionMessage{
			{
				Role: openai.ChatMessageRoleSystem,
				Content: "You are a Go developer. You use the hexagonal architecture pattern.\n" +
					"Do not write any accompanying explanations, only provide the code without quoting it.\n" +
					"Only answer with a single message.\n" +
					"You are writing a Go package named " + pkg + " that implements the interface " + spec + ".\n" +
					"The user will provide a specification of the interace and you will respond with " +
					"a Go package that implements the interface.\n" +
					"Do not include the original code in the response.\n" +
					"Your code must compile alongside the user's code as part of the same package.\n",
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: string(input),
			},
		},
	})
	if err != nil {
		return err
	}
	// output, err := format.Source([]byte(resp.Choices[0].Message.Content))
	// if err != nil {
	// 	return err
	// }
	return os.WriteFile(out, []byte(resp.Choices[0].Message.Content), 0644)
}
