package llm

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewOpenAIClient_DefaultBaseURL(t *testing.T) {
	client := NewOpenAIClient("", "test-key", "test-model")
	if client.BaseURL != "https://api.openai.com/v1" {
		t.Errorf("Expected default base URL, got %s", client.BaseURL)
	}
}

func TestNewOpenAIClient_CustomBaseURL(t *testing.T) {
	client := NewOpenAIClient("https://custom.api.com/v1", "key", "model")
	if client.BaseURL != "https://custom.api.com/v1" {
		t.Errorf("Expected custom base URL, got %s", client.BaseURL)
	}
}

func TestChatWithTools_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST, got %s", r.Method)
		}
		if r.Header.Get("Authorization") != "Bearer test-key" {
			t.Errorf("Expected Bearer token, got %s", r.Header.Get("Authorization"))
		}

		var req chatRequest
		json.NewDecoder(r.Body).Decode(&req)
		if req.Model != "test-model" {
			t.Errorf("Expected model test-model, got %s", req.Model)
		}

		resp := chatResponse{
			Choices: []struct {
				Message Message `json:"message"`
			}{
				{Message: Message{Role: "assistant", Content: "Hello from LLM"}},
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewOpenAIClient(server.URL, "test-key", "test-model")
	msg, err := client.ChatWithTools([]Message{{Role: "user", Content: "Hi"}}, nil)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if msg.Content != "Hello from LLM" {
		t.Errorf("Expected 'Hello from LLM', got %s", msg.Content)
	}
}

func TestChatWithTools_APIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": map[string]string{"message": "internal server error"},
		})
	}))
	defer server.Close()

	client := NewOpenAIClient(server.URL, "test-key", "test-model")
	_, err := client.ChatWithTools([]Message{{Role: "user", Content: "Hi"}}, nil)
	if err == nil {
		t.Error("Expected error for 500 response")
	}
}

func TestChatWithTools_EmptyChoices(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := chatResponse{Choices: []struct {
			Message Message `json:"message"`
		}{}}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewOpenAIClient(server.URL, "test-key", "test-model")
	_, err := client.ChatWithTools([]Message{{Role: "user", Content: "Hi"}}, nil)
	if err == nil {
		t.Error("Expected error for empty choices")
	}
}

func TestChatCompletion_WrapsChatWithTools(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := chatResponse{
			Choices: []struct {
				Message Message `json:"message"`
			}{
				{Message: Message{Role: "assistant", Content: "test response"}},
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewOpenAIClient(server.URL, "key", "model")
	content, err := client.ChatCompletion([]Message{{Role: "user", Content: "Hi"}})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if content != "test response" {
		t.Errorf("Expected 'test response', got %s", content)
	}
}

func TestChatWithTools_ConnectionError(t *testing.T) {
	client := NewOpenAIClient("http://localhost:1", "key", "model")
	_, err := client.ChatWithTools([]Message{{Role: "user", Content: "Hi"}}, nil)
	if err == nil {
		t.Error("Expected error for connection failure")
	}
}
