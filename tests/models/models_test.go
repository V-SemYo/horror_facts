package models_test

import (
	"horror_facts/models"
	"testing"
)

func TestMovieModel(t *testing.T) {
	movie := models.Movie{
		Title:    "Сияние",
		Year:     1980,
		Category: "зарубежный",
	}
	if movie.Title != "Сияние" {
		t.Errorf("Title mismatch: got %s", movie.Title)
	}
	if movie.Year != 1980 {
		t.Errorf("Year mismatch: got %d", movie.Year)
	}
	if movie.Category != "зарубежный" {
		t.Errorf("Category mismatch: got %s", movie.Category)
	}
}
