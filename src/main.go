package main

import (
	"flag"
	"log"
	"net"
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
		listener, err := net.Listen("tcp", ":8080")
		if err != nil {
			log.Fatalf("Failed to start server: %v", err)
			return
		}

		defer listener.Close()
		log.Println("server is running on port 8080...")

		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Printf("failed to accept connection: %v", err)
				continue
			}
			go handleConnection(conn, model, *cfg)
		}
	}
}
