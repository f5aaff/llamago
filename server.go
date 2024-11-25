package main


import (
	"bufio"
	"fmt"
	"log"
	"net"
	"encoding/json"
	llama "github.com/go-skynet/go-llama.cpp"
)

func handleConnection(conn net.Conn, model *llama.LLama, cfg Config) {
	defer conn.Close()
	history := []string{cfg.InitPrompt}

	reader := bufio.NewReader(conn)

	for {
		message, err := reader.ReadString('\n')

		if err != nil {
			log.Printf("error reading from client: %v", err)
		}

		var req Request
		if err := json.Unmarshal([]byte(message), &req); err != nil {
			log.Printf("error parsing JSON: %v", err)
			resp := Response{Error: "Invalid JSON format"}
			sendResponse(conn, resp)
			continue
		}

		history = append(history, "User: "+req.Message)

		respText, err := genResponse(&history, model, cfg)
		resp := Response{Response: respText}

		if err != nil {
			resp.Error = fmt.Sprintf("MOdel error: %v", err)
		}

		sendResponse(conn, resp)
	}
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
