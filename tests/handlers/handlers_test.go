package handler_test

import (
	"horror_facts/handlers"
	"horror_facts/models"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHomeHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handlers.HomeHandler(rr, req)

	t.Logf("Status: %d", rr.Code)
	t.Logf("Body length: %d bytes", rr.Body.Len())

	if status := rr.Code; status != http.StatusOK {
		t.Logf("Full response body: %s", rr.Body.String())
		t.Errorf("Expected status 200, got %d", status)
	}

	body := rr.Body.String()
	if len(body) == 0 {
		t.Error("Empty response body")
	}
}

func TestSearchHandler(t *testing.T) {
	testMovies := map[string]models.Movie{
		"сияние": {
			Title:    "Сияние",
			Year:     1980,
			About:    "Тестовое описание",
			Facts:    "Тестовые факты",
			Category: "зарубежный",
		},
	}

	tests := []struct {
		name     string
		query    string
		wantCode int
		wantBody string
	}{
		{
			name:     "Existing movie",
			query:    "сияние",
			wantCode: http.StatusOK,
			wantBody: "Сияние",
		},
		{
			name:     "Non-existing movie",
			query:    "несуществующий",
			wantCode: http.StatusOK,
			wantBody: "не найден",
		},
		{
			name:     "Empty query",
			query:    "",
			wantCode: http.StatusOK,
			wantBody: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/search?q="+tt.query, nil)
			rr := httptest.NewRecorder()

			// ПРАВИЛЬНЫЙ ВЫЗОВ:
			handler := handlers.SearchHandler(testMovies)
			handler.ServeHTTP(rr, req)

			if rr.Code != tt.wantCode {
				t.Errorf("Status: got %d, want %d", rr.Code, tt.wantCode)
			}

			if tt.wantBody != "" && !strings.Contains(rr.Body.String(), tt.wantBody) {
				t.Errorf("Body should contain %q, got: %s",
					tt.wantBody, rr.Body.String())
			}
		})
	}
}
