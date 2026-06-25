package handlers

import (
	"encoding/json"
	"horror_facts/internal/ai"
	"horror_facts/models"
	"html/template"
	"net/http"
	"strings"
)

func SearchHandler(movies map[string]models.Movie) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var query string
		if r.Method == "POST" {
			query = r.FormValue("movie")
		} else {
			query = r.URL.Query().Get("q")
		}

		searchKey := strings.ToLower(strings.TrimSpace(query))
		movie, found := movies[searchKey]

		if !found {
			moviesJSON, _ := json.Marshal(movies)
			aiResponse, err := ai.SearchMovie(query, string(moviesJSON))
			if err == nil {
				var aiMovie models.Movie
				if json.Unmarshal([]byte(aiResponse), &aiMovie) == nil {
					tmpl := template.Must(template.ParseFiles("templates/base.html", "templates/search.html"))
					tmpl.ExecuteTemplate(w, "base.html", map[string]interface{}{
						"Title":      aiMovie.Title,
						"MovieTitle": aiMovie.Title,
						"Year":       aiMovie.Year,
						"Category":   aiMovie.Category,
						"About":      aiMovie.About,
						"Facts":      aiMovie.Facts,
						"Username":   GetUserEmail(r),
					})
					return
				}
			}
			tmpl := template.Must(template.ParseFiles("templates/base.html", "templates/notfound.html"))
			tmpl.ExecuteTemplate(w, "base.html", map[string]interface{}{
				"Title": "Not Found",
				"Query": query,
			})
			return
		}

		tmpl := template.Must(template.ParseFiles("templates/base.html", "templates/search.html"))
		tmpl.ExecuteTemplate(w, "base.html", map[string]any{
			"Title":      movie.Title,
			"MovieTitle": movie.Title,
			"Year":       movie.Year,
			"Category":   movie.Category,
			"About":      movie.About,
			"Facts":      movie.Facts,
			"Username":   GetUserEmail(r),
		})
	}
}
