package handlers

import (
	"horror_facts/models"
	"html/template"
	"net/http"
	"strings"
)

func AllMoviesHandler(movies map[string]models.Movie) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var russianMovies, foreignMovies []struct {
			Key   string
			Movie models.Movie
		}
		for key, movie := range movies {
			if strings.EqualFold(movie.Category, "русский") {
				russianMovies = append(russianMovies, struct {
					Key   string
					Movie models.Movie
				}{key, movie})
			} else {
				foreignMovies = append(foreignMovies, struct {
					Key   string
					Movie models.Movie
				}{key, movie})
			}
		}

		tmpl := template.Must(template.ParseFiles("templates/base.html", "templates/movies.html"))
		tmpl.ExecuteTemplate(w, "base.html", map[string]any{
			"Title":         "Archive",
			"RussianMovies": russianMovies,
			"ForeignMovies": foreignMovies,
			"Username":      GetUserEmail(r),
		})
	}
}
