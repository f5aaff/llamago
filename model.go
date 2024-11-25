package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	llama "github.com/go-skynet/go-llama.cpp"
)

// takes input in basic fashion from terminal.
func takeInput(history *[]string) error {
	fmt.Print("msg: ")
	reader := bufio.NewReader(os.Stdin)
	userIn, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	*history = append(*history, "User: "+userIn)
	return nil
}

// takes the history, feeds it in as a prompt, and returns the generated response.
func genResponse(history *[]string, model *llama.LLama, cfg Config) (string, error) {
	prompt := strings.Join(*history, "\n") + "\n:Assistant:"
	resp, err := model.Predict(prompt, llama.SetTemperature(0.7), llama.SetTopK(50), llama.SetTokens(cfg.MaxTokens), llama.SetStopWords("User:"))
	if err != nil {
		return "", err
	}

	*history = append(*history, "Assistant: "+resp)
	fmt.Println("Assistant: " + resp)
	return resp, nil
}

// for loop to continuously take input and spit out responses.
func initConversation(history *[]string, model *llama.LLama, cfg Config) error {
	println("starting conversation...")
	for {
		err := takeInput(history)
		if err != nil {
			println("TAKEINPUT ERROR")
			return err
		}
		_, err = genResponse(history, model, cfg)
		if err != nil {
			println("GENRESPONSE ERROR")
			return err
		}
	}
}
