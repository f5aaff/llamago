package main


// Simple struct to represent config opts. from file.
type Config struct {
	ModelPath         string `json:"model_path"`
	InitPrompt        string `json:"init_prompt"`
	MaxTokens         int    `json:"max_tokens"`
	ModelContextLimit int    `json:"model_context_limit"`
}

type Request struct {
	Message string `json:"message"`
}

type Response struct {
	Response string `json:"response"`
	Error    string `json:"error,omitempty"`
}
