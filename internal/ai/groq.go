package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type GroqResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func GenerateCommitMessages(diff string) ([]string, error) {
	apiKey := os.Getenv("GROQ_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("missing GROQ_API_KEY")
	}

	payload := map[string]interface{}{
		"model": "meta-llama/llama-4-scout-17b-16e-instruct",
		"messages": []map[string]string{
			{
				"role":    "system",
				"content": "You generate SHORT, clean conventional commit messages. Always return exactly 5 options.",
			},
			{
				"role":    "user",
				"content": fmt.Sprintf("Generate 5 conventional commit messages for the following git diff:\n\n%s", diff),
			},
		},
		"temperature": 0.3,
		"max_tokens":  200,
	}

	bodyBytes, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST",
		"https://api.groq.com/openai/v1/chat/completions",
		bytes.NewBuffer(bodyBytes),
	)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP error: %v", err)
	}
	defer resp.Body.Close()

	raw, _ := io.ReadAll(resp.Body)

	// Debug mode
	if os.Getenv("CODESCRIBE_DEBUG") == "1" {
		fmt.Println("RAW GROQ RESPONSE:", string(raw))
	}

	var out GroqChatResponse
	if err := json.Unmarshal(raw, &out); err != nil {
		return nil, fmt.Errorf("JSON parse error: %v\nRAW: %s", err, string(raw))
	}

	// 🔥 FIX: ensure Groq returned choices
	if len(out.Choices) == 0 || out.Choices[0].Message.Content == "" {
		return nil, fmt.Errorf("empty Groq response: %s", string(raw))
	}

	// Split commit options by line
	content := out.Choices[0].Message.Content
	lines := strings.Split(content, "\n")

	var msgs []string
	for _, line := range lines {
		trim := strings.TrimSpace(line)
		if trim != "" {
			msgs = append(msgs, trim)
		}
	}

	return msgs, nil
}
