package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"

	llama "github.com/go-skynet/go-llama.cpp"
)

func handleConnection(conn net.Conn, model *llama.LLama, cfg *Config) {
	defer conn.Close()
	history := []string{cfg.InitPrompt}

	reader := bufio.NewReader(conn)

		message,err := reader.ReadString('\n')
		fmt.Printf("MESSAGE:%s\n",string(message))
		if err != nil {
			log.Printf("error reading from client: %v", err)
			return
		}

		var req Request
		if err := json.Unmarshal([]byte(message), &req); err != nil {
			log.Printf("error parsing JSON: %v", err)
			resp := Response{Error: "Invalid JSON format"}
			sendResponse(conn, resp)
			return
		}

		history = append(history, "User: "+req.Message)

		respText, err := genResponse(&history, model, *cfg)
		resp := Response{Response: respText}

		if err != nil {
			resp.Error = fmt.Sprintf("Model error: %v", err)
		}

		sendResponse(conn, resp)
}

func sendResponse(conn net.Conn, resp Response) {
	respJson, err := json.Marshal(resp)
	if err != nil {
		log.Printf("error encoding JSON response: %v", err)
		return
	}

	_, err = conn.Write(append(respJson, '\n'))
	if err != nil {
		log.Printf("error writing to client: %v", err)
	}
}

func StartServer(cfg *Config, model *llama.LLama) error {
	listener, err := net.Listen("tcp", cfg.ListenAddress)
	if err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}

	defer listener.Close()

	log.Printf("Server is running on %s...", cfg.ListenAddress)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			continue
		}
		go handleConnection(conn, model, cfg)
	}
}
