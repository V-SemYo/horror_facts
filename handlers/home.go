package handlers

import (
	"html/template"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

// HomeHandler отображает главную страницу
func HomeHandler(w http.ResponseWriter, r *http.Request) {

	tmpl := template.Must(template.ParseFiles("templates/base.html", "templates/home.html"))
	tmpl.ExecuteTemplate(w, "base.html", map[string]any{
		"Title":    "Home",
		"Username": GetUserEmail(r),
	})
}

// GetUserEmail достаёт email пользователя из JWT токена в cookie
func GetUserEmail(r *http.Request) string {
	cookie, err := r.Cookie("token")
	if err != nil {
		return ""
	}

	token, _ := jwt.Parse(cookie.Value, func(t *jwt.Token) (any, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if email, ok := claims["email"].(string); ok {
			return email
		}
	}

	return ""
}
