package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"horror_facts/internal/config"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type MovieJSON struct {
	Key      string `json:"key"`
	Title    string `json:"title"`
	Year     int    `json:"year"`
	About    string `json:"about"`
	Facts    string `json:"facts"`
	Category string `json:"category"`
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ .env not found")
	}

	cfg := config.Load()

	db, err := sql.Open("postgres", cfg.DataBase.DSN())
	if err != nil {
		log.Fatalf("❌ DB open error: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("❌ DB ping error: %v", err)
	}

	log.Println("✅ Connected to DB")

	file, err := os.ReadFile("data/movies.json")
	if err != nil {
		log.Fatalf("❌ Read file error: %v", err)
	}

	var movies []MovieJSON
	if err := json.Unmarshal(file, &movies); err != nil {
		log.Fatalf("❌ JSON parse error: %v", err)
	}

	log.Printf("📄 Found %d movies in JSON", len(movies))

	inserted := 0
	for _, m := range movies {
		_, err := db.Exec(
			`INSERT INTO movies (key, title, year, about, facts, category) 
             VALUES ($1, $2, $3, $4, $5, $6) 
             ON CONFLICT (key) DO NOTHING`,
			m.Key, m.Title, m.Year, m.About, m.Facts, m.Category,
		)
		if err != nil {
			log.Printf("⚠️ Error inserting '%s': %v", m.Title, err)
			continue
		}
		inserted++
	}

	log.Printf("✅ Inserted %d movies into DB", inserted)
	fmt.Println("Done!")
}
