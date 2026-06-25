package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

// Структура для запроса к API
type ChatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

// Структура сообщения для API
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// Структура для ответа от API
type ChatResponse struct {
	Choices []Choice `json:"choices"`
}

// Выбор нашего ответа сообщения от API
type Choice struct {
	Message Message `json:"message"`
}

// SearchMovie отправляет запрос в DeepSeek через OpenRouter
// query - что пользователь ввёл
// moviesContext - 65 фильмов из БД в виде JSON строки
func SearchMovie(query string, moviesContext string) (string, error) {
	apiKey := os.Getenv("OPENROUTER_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("OPENROUTER_API_KEY not set")
	}

	prompt := fmt.Sprintf(`You are a horror movie database. Given the user's search query и на русском, find the most relevant movie from this list:

%s

If the query matches a movie in the list, return ONLY a JSON object with these exact fields:
{
  "title": "Movie Title",
  "year": 1980,
  "about": "Brief description",
  "facts": "Real facts behind the movie",
  "category": "русский or зарубежный"
}

If the query does NOT match any movie, generate a plausible answer in the same JSON format about a famous horror movie.

User query: %s

Return ONLY the JSON object, nothing else.`, moviesContext, query)

	reqBody := ChatRequest{
		Model: "deepseek/deepseek-chat",
		Messages: []Message{
			{Role: "user", Content: prompt},
		},
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", "https://openrouter.ai/api/v1/chat/completions", bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var chatResp ChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&chatResp); err != nil {
		return "", err
	}

	if len(chatResp.Choices) == 0 {
		return "", fmt.Errorf("no response from AI")
	}

	content := chatResp.Choices[0].Message.Content

	// Очищаем ответ ИИ от Markdown-обёртки ```json ... ```
	content = strings.TrimSpace(content)
	content = strings.TrimPrefix(content, "```json")
	content = strings.TrimPrefix(content, "```")
	content = strings.TrimSuffix(content, "```")
	content = strings.TrimSpace(content)

	return content, nil
}
