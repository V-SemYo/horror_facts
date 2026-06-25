package handlers

import (
	"horror_facts/models"
	"html/template"
	"net/http"
	"strings"
)

func FilterHandler(movies map[string]models.Movie) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		category := r.URL.Query().Get("category")
		if category == "" {
			http.Redirect(w, r, "/movies", http.StatusSeeOther)
			return
		}

		var filtered []struct {
			Key   string
			Movie models.Movie
		}
		for key, movie := range movies {
			if strings.EqualFold(movie.Category, category) {
				filtered = append(filtered, struct {
					Key   string
					Movie models.Movie
				}{key, movie})
			}
		}

		categoryTitle := "Foreign Horror"
		if strings.EqualFold(category, "русский") {
			categoryTitle = "Russian Horror"
		}

		tmpl := template.Must(template.ParseFiles("templates/base.html", "templates/filter.html"))
		tmpl.ExecuteTemplate(w, "base.html", map[string]any{
			"Title":         categoryTitle,
			"CategoryTitle": categoryTitle,
			"Count":         len(filtered),
			"Movies":        filtered,
			"Username":      GetUserEmail(r),
		})
	}
}
