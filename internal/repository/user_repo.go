package repository

import (
	"database/sql"
	"horror_facts/models"
)

// UserRepo работает с таблицей users
type UserRepo struct {
	db *sql.DB
}

// NewUserRepo создаёт новый пользовательский репозиторий
func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

// CreateUser создаёт нового пользователя и возвращает его
func (ur *UserRepo) CreateUser(username, email, passwordHash string) (*models.User, error) {
	user := &models.User{}

	err := ur.db.QueryRow(`INSERT INTO users (username, email, password_hash)
	VALUES ($1, $2, $3) 
         RETURNING id, username, email, created_at
	`, username, email, passwordHash).Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt)

	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetByEmail ищет пользователя по email
func (ur *UserRepo) GetByEmail(email string) (*models.User, error) {
	user := &models.User{}

	err := ur.db.QueryRow(`SELECT id, username, email, password_hash, created_at
	FROM users WHERE email = $1`, email).Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.CreatedAt)

	if err != nil {
		return nil, err
	}

	return user, nil
}
