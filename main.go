package main

import (
	"encoding/json"
	"fmt"
	"horror_facts/handlers"
	"horror_facts/internal/config"
	"horror_facts/internal/database"
	"horror_facts/internal/middleware"
	"horror_facts/internal/repository"
	"horror_facts/models"
	"html/template"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/joho/godotenv"
)

var moviesCache map[string]models.Movie

var cashMutex sync.RWMutex

type MovieJSON struct {
	Key      string `json:"key"`
	Title    string `json:"title"`
	Year     int    `json:"year"`
	About    string `json:"about"`
	Facts    string `json:"facts"`
	Category string `json:"category"`
}

func loadMovies() (map[string]models.Movie, error) {
	file, err := os.Open("data/movies.json")
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	var moviesList []MovieJSON

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&moviesList); err != nil {
		return nil, fmt.Errorf("JSON parsing error: %v", err)
	}

	if len(moviesList) == 0 {
		return nil, fmt.Errorf("file is empty")
	}

	movies := make(map[string]models.Movie)
	validCount := 0

	for _, m := range moviesList {
		if m.Key == "" || m.Title == "" {
			log.Printf("Skipped movie with empty key or title")
			continue
		}
		movies[m.Key] = models.Movie{
			Title:    m.Title,
			Year:     m.Year,
			About:    m.About,
			Facts:    m.Facts,
			Category: m.Category,
		}
		validCount++
	}

	if validCount == 0 {
		return nil, fmt.Errorf("no valid entries in file")
	}

	log.Printf("Loaded %d movies from JSON", validCount)
	return movies, nil
}

func getMovieFromCashe(key string) (models.Movie, bool) {
	cashMutex.RLock()
	defer cashMutex.RUnlock()
	movie, found := moviesCache[key]
	return movie, found
}

func getAllMoviesFromCashe() map[string]models.Movie {
	cashMutex.RLock()
	defer cashMutex.RUnlock()
	copyMap := make(map[string]models.Movie)
	for k, v := range moviesCache {
		copyMap[k] = v
	}
	return copyMap
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ .env file not found, using system environment")
	}

	cfg := config.Load()

	log.Printf("DSN: %s", cfg.DataBase.DSN())

	db, err := database.NewPostgreDB(cfg.DataBase)
	if err != nil {
		log.Fatalf("❌ Database error: %v", err)
	}

	userRepo := repository.NewUserRepo(db)

	defer db.Close()

	log.Println("Database connected successfully!")

	movies, err := loadMovies()
	if err != nil {
		log.Fatalf("Database loading error: %v", err)
	}

	cashMutex.Lock()
	moviesCache = movies
	cashMutex.Unlock()

	log.Printf("Database: %d movies", len(movies))
	log.Printf("Cache loaded: %d movies in memory", len(moviesCache))

	mux := http.NewServeMux()

	authHandler := handlers.NewAuthHandler(userRepo)

	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	mux.HandleFunc("/", handlers.HomeHandler)

	mux.HandleFunc("/search", handlers.SearchHandler(movies))

	mux.HandleFunc("/movies", handlers.AllMoviesHandler(movies))

	mux.HandleFunc("/filter", handlers.FilterHandler(movies))

	mux.HandleFunc("/map", handlers.MapHandler)

	mux.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			tmpl := template.Must(template.ParseFiles("templates/base.html", "templates/register.html"))
			tmpl.ExecuteTemplate(w, "base.html", map[string]any{
				"Title": "Register"})
			return
		}
		authHandler.Register(w, r)
	})

	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			tmpl := template.Must(template.ParseFiles("templates/base.html", "templates/login.html"))
			tmpl.ExecuteTemplate(w, "base.html", map[string]any{
				"Title": "Login"})
			return
		}
		authHandler.Login(w, r)
	})

	mux.HandleFunc("/profile", middleware.Auth(func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value(middleware.UserIDKey).(int)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"user_id": userID,
			"message": "This is a protected route",
		})
	}))

	mux.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name:   "token",
			Value:  "",
			Path:   "/",
			MaxAge: -1,
		})
		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

	log.Println("Starting horror website...")
	log.Printf("Movies in database: %d", len(movies))

	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
		ErrorLog:     log.New(os.Stderr, "SERVER: ", log.LstdFlags),
	}

	log.Printf("🌐 Server running on port %s", server.Addr)
	log.Println("VHS Horror site is ready...")

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("❌ Server error: %v", err)
	}
}
