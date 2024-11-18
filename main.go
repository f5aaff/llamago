package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	llama "github.com/go-skynet/go-llama.cpp"
)

// Simple struct to represent config opts. from file.
type Config struct {
	ModelPath  string `json:"model_path"`
	InitPrompt string `json:"init_prompt"`
	MaxTokens int `json:"max_tokens"`
	ModelContextLimit int `json:"model_context_limit"`
}

func load_config(cfg *Config) error {
	f, err := os.Open("config.json")
	if err != nil {
		return err
	}
	defer f.Close()

	b, err := io.ReadAll(f)
	if err != nil {
		return err
	}
	err = json.Unmarshal(b, &cfg)
	if err != nil {
		return err
	}

	return nil
}

// takes input in basic fashion from terminal.
func takeInput(history *[]string) error {
	fmt.Print("msg: ")
	reader := bufio.NewReader(os.Stdin)
	userIn,err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	*history = append(*history, "User: "+userIn)
	return nil
}

// takes the history, feeds it in as a prompt, and returns the generated response.
func genResponse(history *[]string, model *llama.LLama,cfg Config) error {
	prompt := strings.Join(*history, "\n") + "\n:Assistant:"
	resp, err := model.Predict(prompt, llama.SetTemperature(0.7), llama.SetTopK(50), llama.SetTokens(cfg.MaxTokens),llama.SetStopWords("User:"))
	if err != nil {
		return err
	}

	*history = append(*history, "Assistant: "+resp)
	fmt.Println("Assistant: "+resp)
	return nil
}

// for loop to continuously take input and spit out responses.
func initConversation(history *[]string, model *llama.LLama,cfg Config) error {
	println("starting conversation...")
	for {
		err := takeInput(history)
		if err != nil {
			println("TAKEINPUT ERROR")
			return err
		}
		err = genResponse(history, model,cfg)
		if err != nil {
			println("GENRESPONSE ERROR")
			return err
		}
	}
}
func main() {
	cfg := Config{
		ModelPath:  "./models/TinyLLama-1.1B-Chat-v1.0-GGUF/ggml-model-f16.gguf",
		InitPrompt: "You are a chat bot, intended to answer general queries.",
		MaxTokens: 200,
		ModelContextLimit: 2000,
	}

	err := load_config(&cfg)
	if err != nil {
		log.Printf("error loading config, using defaults...")
	}
	// Load the model
	modelPath := cfg.ModelPath
	model, err := llama.New(modelPath, llama.SetContext(cfg.ModelContextLimit))
	if err != nil {
		log.Fatalf("Failed to load model: %v", err)
	}
	history := make([]string, 1)
	err = initConversation(&history, model,cfg)
	if err != nil {
		log.Printf("crash out, shit's fucked: %s\n", err.Error())
	}
}
