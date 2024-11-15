package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	llama "github.com/go-skynet/go-llama.cpp"
)

type Config struct {
	ModelPath  string `json:"model_path"`
	InitPrompt string `json:"init_prompt"`
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

func main() {
	cfg := Config{
		ModelPath: "./models/TinyLLama-1.1B-Chat-v1.0-GGUF/ggml-model-f16.gguf",
		InitPrompt: "You are a chat bot, intended to answer general queries.",
	}

	err := load_config(&cfg)
	if err != nil {
		log.Printf("error loading config, using defaults...")
	}
	// Load the model
	modelPath := cfg.ModelPath
	model, err := llama.New(modelPath, llama.SetContext(512))
	if err != nil {
		log.Fatalf("Failed to load model: %v", err)
	}

	// Prompt for the model
	prompt := cfg.InitPrompt
	result, err := model.Predict(prompt, llama.SetThreads(4))
	if err != nil {
		log.Fatalf("Prediction error: %v", err)
	}
	fmt.Printf("Model Response: %s\n", result)
}
