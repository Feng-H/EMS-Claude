package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type LLMClient interface {
	ChatCompletion(messages []Message) (string, error)
	ChatWithTools(messages []Message, tools []Tool) (*Message, error)
}

type Message struct {
	Role       string     `json:"role"`
	Content    string     `json:"content"`
	ToolCalls  []ToolCall `json:"tool_calls,omitempty"`
	ToolCallID string     `json:"tool_call_id,omitempty"`
}

type ToolCall struct {
	ID       string `json:"id"`
	Type     string `json:"type"`
	Function struct {
		Name      string `json:"name"`
		Arguments string `json:"arguments"`
	} `json:"function"`
}

type Tool struct {
	Type     string `json:"type"`
	Function struct {
		Name        string      `json:"name"`
		Description string      `json:"description"`
		Parameters  interface{} `json:"parameters"`
	} `json:"function"`
}

type OpenAIClient struct {
	BaseURL string
	APIKey  string
	Model   string
}

func NewOpenAIClient(baseURL, apiKey, model string) *OpenAIClient {
	if baseURL == "" {
		baseURL = "https://api.openai.com/v1"
	}
	return &OpenAIClient{
		BaseURL: baseURL,
		APIKey:  apiKey,
		Model:   model,
	}
}

type chatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Tools    []Tool    `json:"tools,omitempty"`
}

type chatResponse struct {
	Choices []struct {
		Message Message `json:"message"`
	} `json:"choices"`
}

func (c *OpenAIClient) ChatCompletion(messages []Message) (string, error) {
	resp, err := c.ChatWithTools(messages, nil)
	if err != nil {
		return "", err
	}
	return resp.Content, nil
}

func (c *OpenAIClient) ChatWithTools(messages []Message, tools []Tool) (*Message, error) {
	reqBody := chatRequest{
		Model:    c.Model,
		Messages: messages,
		Tools:    tools,
	}
	
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", c.BaseURL+"/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.APIKey)

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp struct {
			Error struct {
				Message string `json:"message"`
			} `json:"error"`
		}
		json.NewDecoder(resp.Body).Decode(&errResp)
		return nil, fmt.Errorf("LLM API error (status %d): %s", resp.StatusCode, errResp.Error.Message)
	}

	var result chatResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if len(result.Choices) > 0 {
		return &result.Choices[0].Message, nil
	}

	return nil, fmt.Errorf("no response from LLM")
}
