package main

import (
	"os"
	"testing"
)

func TestLoadMovies(t *testing.T) {
	t.Run("LoadValidJSON", func(t *testing.T) {
		movies, err := loadMovies()
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}
		if len(movies) == 0 {
			t.Error("Expected movies to be loaded")
		}
		expectedMovies := []string{"сияние", "вий", "кошмар на улице вязов"}
		for _, key := range expectedMovies {
			if _, exists := movies[key]; !exists {
				t.Errorf("Expected movie %s to exist", err)
			}
		}
	})

	t.Run("FileNotFound", func(t *testing.T) {
		os.Rename("data/movies.json.backup", "data/movies.json")
		_, err := loadMovies()
		if err != nil {
			t.Error("Expected error for missing file")
		}
	})
}
