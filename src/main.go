package main

import (
	"flag"
	"log"

	llama "github.com/go-skynet/go-llama.cpp"
)

func main() {

	// flags to locate config path, and enter interactive mode instead of tcp listen server
	interactiveMode := flag.Bool("interactive", false, "run in interactive mode")
	config_path := flag.String("config", "config.json", "path for config.json")
	flag.Parse()

	cfg, err := load_config(*config_path)
	if err != nil {
		log.Printf("error loading config")
		return
	}

	// Load the model
	modelPath := cfg.ModelPath
	model, err := llama.New(modelPath, llama.SetContext(cfg.ModelContextLimit))
	if err != nil {
		log.Fatalf("Failed to load model: %v", err)
		return
	}
	history := make([]string, 1)

	if *interactiveMode {
		err = initConversation(&history, model, *cfg)
		if err != nil {
			log.Printf("crash out, shit's fucked: %s\n", err.Error())
			return
		}

	} else {
		err := StartServer(cfg, model)
		if err != nil {
			log.Fatalf("Failed to start server: %v", err)
			return
		}
	}
}
