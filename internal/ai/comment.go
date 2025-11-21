package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type commentResp struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func AddCommentsToFile(path string) error {
	code, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	commented, err := requestAICommenting(string(code))
	if err != nil {
		// fallback
		fallback := fmt.Sprintf("// AI failed to process. Original code:\n\n%s", string(code))
		return ioutil.WriteFile(path, []byte(fallback), 0644)
	}

	return ioutil.WriteFile(path, []byte(commented), 0644)
}

func requestAICommenting(code string) (string, error) {
	apiKey := os.Getenv("GROQ_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("missing GROQ_API_KEY")
	}

	body := map[string]interface{}{
		"model": "llama3-70b-8192",
		"messages": []map[string]string{
			{
				"role": "user",
				"content": fmt.Sprintf(
					"Add detailed inline comments to the following code. Detect its language automatically. Preserve formatting. Return only the commented code.\n\n--- CODE START ---\n%s\n--- CODE END ---",
					code,
				),
			},
		},
	}

	j, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "https://api.groq.com/openai/v1/chat/completions", bytes.NewBuffer(j))

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var out commentResp
	json.NewDecoder(resp.Body).Decode(&out)

	if len(out.Choices) == 0 {
		return "", fmt.Errorf("empty response")
	}

	return out.Choices[0].Message.Content, nil
}
