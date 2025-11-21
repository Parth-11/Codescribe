package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
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

	body := map[string]interface{}{
		"model": "llama3-70b-8192",
		"messages": []map[string]string{
			{
				"role":    "user",
				"content": fmt.Sprintf("Generate 5 conventional commit messages based on this diff:\n\n%s", diff),
			},
		},
	}

	jsonData, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "https://api.groq.com/openai/v1/chat/completions", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var groqResponse GroqResponse
	json.NewDecoder(resp.Body).Decode(&groqResponse)

	if len(groqResponse.Choices) == 0 {
		return nil, fmt.Errorf("empty Groq response")
	}

	raw := groqResponse.Choices[0].Message.Content

	// split into list by newlines
	lines := bytes.Split([]byte(raw), []byte("\n"))

	var msgs []string
	for _, l := range lines {
		trim := string(bytes.TrimSpace(l))
		if trim != "" {
			msgs = append(msgs, trim)
		}
	}

	return msgs, nil
}
