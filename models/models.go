package models

import "time"

type Movie struct {
	Title    string `json:"title"`
	Year     int    `json:"year"`
	About    string `json:"about"`
	Facts    string `json:"facts"`
	Category string `json:"category"`
}

type User struct {
	ID           int       `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
}
