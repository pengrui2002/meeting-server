package repository

import (
    "github.com/jmoiron/sqlx"
    "meeting-server/internal/model"
)

type UserRepository struct {
    db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
    return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *model.User) error {
    query := `INSERT INTO users (id, username, password_hash, display_name) VALUES ($1, $2, $3, $4)`
    _, err := r.db.Exec(query, user.ID, user.Username, user.PasswordHash, user.DisplayName)
    return err
}

func (r *UserRepository) GetByUsername(username string) (*model.User, error) {
    var user model.User
    query := `SELECT * FROM users WHERE username = $1`
    err := r.db.Get(&user, query, username)
    if err != nil {
        return nil, err
    }
    return &user, nil
}

func (r *UserRepository) GetByID(id string) (*model.User, error) {
    var user model.User
    query := `SELECT * FROM users WHERE id = $1`
    err := r.db.Get(&user, query, id)
    if err != nil {
        return nil, err
    }
    return &user, nil
}