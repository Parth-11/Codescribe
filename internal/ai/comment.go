package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

const maxChunkSize = 12000 // safe chunk size for Groq

// Generic Groq response struct
type GroqChatResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

// ----------------------------
// MAIN FUNC – COMMENT A FILE
// ----------------------------

func AddCommentsToFile(path string) error {
	codeBytes, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	code := string(codeBytes)

	// Split into small chunks to avoid token overflow
	chunks := splitIntoChunks(code, maxChunkSize)

	var commentedBuilder strings.Builder

	for i, chunk := range chunks {
		commented, err := requestAICommenting(chunk)
		if err != nil {
			fmt.Printf("⚠️ AI failed for chunk %d in %s: %v\n", i+1, path, err)
			commentedBuilder.WriteString("\n// AI FAILED TO COMMENT THIS SECTION. ORIGINAL CODE:\n")
			commentedBuilder.WriteString(chunk)
			commentedBuilder.WriteString("\n\n")
			continue
		}

		commentedBuilder.WriteString(commented)
		commentedBuilder.WriteString("\n\n")
	}

	// Write commented output
	return os.WriteFile(path, []byte(commentedBuilder.String()), 0644)
}

// ----------------------------
// CHUNKING SYSTEM
// ----------------------------

func splitIntoChunks(text string, size int) []string {
	var chunks []string
	runes := []rune(text)

	for len(runes) > size {
		chunks = append(chunks, string(runes[:size]))
		runes = runes[size:]
	}
	chunks = append(chunks, string(runes))
	return chunks
}

// ----------------------------
// REQUEST AI COMMENTING (with retries)
// ----------------------------

func requestAICommenting(code string) (string, error) {
	apiKey := os.Getenv("GROQ_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("missing GROQ_API_KEY")
	}

	for attempt := 1; attempt <= 3; attempt++ {

		respText, err := callGroqAPI(code)

		if err == nil && strings.TrimSpace(respText) != "" {
			return respText, nil
		}

		fmt.Printf("🔁 Retry %d/3 due to error: %v\n", attempt, err)
		time.Sleep(time.Duration(attempt) * time.Second)
	}

	return "", fmt.Errorf("groq API failed after 3 retries")
}

// ----------------------------
// ACTUAL GROQ API CALL
// ----------------------------

func callGroqAPI(code string) (string, error) {

	// ✨ STRICT SYSTEM PROMPT to prevent rewriting code ✨
	systemPrompt := `
You are a professional code annotator.

RULES YOU MUST FOLLOW STRICTLY:
- DO NOT rewrite the code.
- DO NOT change ANY indentation.
- DO NOT rename ANY variables, functions, or identifiers.
- DO NOT remove ANY code.
- DO NOT reformat the code.
- DO NOT modify blank lines.

YOU MUST:
- Add short, helpful comments ABOVE relevant lines.
- Only add comments where helpful.
- Preserve 100% of the original code EXACTLY.
- Return ONLY the commented code with no explanation.
`

	payload := map[string]interface{}{
		"model": "meta-llama/llama-4-scout-17b-16e-instruct",
		"messages": []map[string]string{
			{
				"role":    "system",
				"content": systemPrompt,
			},
			{
				"role": "user",
				"content": fmt.Sprintf(
					"Please add helpful comments above lines without changing the code:\n\n%s",
					code,
				),
			},
		},
		"temperature": 0.1,
	}

	reqBody, _ := json.Marshal(payload)

	req, _ := http.NewRequest(
		"POST",
		"https://api.groq.com/openai/v1/chat/completions",
		bytes.NewBuffer(reqBody),
	)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+os.Getenv("GROQ_API_KEY"))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("HTTP error: %v", err)
	}
	defer resp.Body.Close()

	raw, _ := io.ReadAll(resp.Body)

	if os.Getenv("CODESCRIBE_DEBUG") == "1" {
		fmt.Println("RAW GROQ RESPONSE:", string(raw))
	}

	var out GroqChatResponse
	if err := json.Unmarshal(raw, &out); err != nil {
		return "", fmt.Errorf("JSON parse error: %v\nRAW:%s", err, string(raw))
	}

	if len(out.Choices) == 0 {
		return "", fmt.Errorf("empty choices returned: %s", string(raw))
	}

	content := out.Choices[0].Message.Content

	if strings.TrimSpace(content) == "" {
		return "", fmt.Errorf("empty message content: %s", string(raw))
	}

	return content, nil
}
