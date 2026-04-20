package repository

import (
	"database/sql"
	"time"

	"github.com/lusiker/clipper/internal/model"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *model.User) error {
	query := `INSERT INTO users (id, username, password, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)`

	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	_, err := r.db.Exec(query, user.ID, user.Username, user.Password, user.CreatedAt, user.UpdatedAt)
	return err
}

func (r *UserRepository) FindByUsername(username string) (*model.User, error) {
	query := `SELECT id, username, password, created_at, updated_at FROM users WHERE username = ?`

	user := &model.User{}
	err := r.db.QueryRow(query, username).Scan(
		&user.ID, &user.Username, &user.Password, &user.CreatedAt, &user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) FindByID(id string) (*model.User, error) {
	query := `SELECT id, username, password, created_at, updated_at FROM users WHERE id = ?`

	user := &model.User{}
	err := r.db.QueryRow(query, id).Scan(
		&user.ID, &user.Username, &user.Password, &user.CreatedAt, &user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}