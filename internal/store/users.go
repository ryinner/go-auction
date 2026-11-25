package store

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/ryinner/go-auction/internal/model"
)

var ErrEmailTaken = errors.New("email already taken")
var ErrUserNotFound = errors.New("user not found")

func (s *Store) CreateUser(ctx context.Context, email string, passwordHash string) (model.User, error) {
	user := model.User{}
	normalizedEmail := strings.ToLower(email)

	err := s.pool.QueryRow(ctx, "INSERT INTO users (email, password_hash) VALUES ($1, $2) RETURNING id, email", normalizedEmail, passwordHash).Scan(&user.ID, &user.Email)

	if err != nil {
		var pgErr *pgconn.PgError
		// Тут as - Другая структура
		if errors.As(err, &pgErr) && pgErr.Code == "23505" && pgErr.ConstraintName == "users_email_key" {
			return model.User{}, ErrEmailTaken
		}

		return model.User{}, fmt.Errorf("creating user error: %w", err)
	}

	return user, nil
}

func (s *Store) GetUserByEmail(ctx context.Context, email string) (model.User, error) {
	normalizedEmail := strings.ToLower(email)

	user := model.User{}

	err := s.pool.QueryRow(ctx, "SELECT id, email, password_hash FROM users WHERE email = $1", normalizedEmail).Scan(&user.ID, &user.Email, &user.PasswordHash)

	if err != nil {
		// Тут is - ErrNoRows = Сравнение ошибки, именно из пакета pbx, если такой же текст = другая ошибка
		if errors.Is(err, pgx.ErrNoRows) {
			return model.User{}, ErrUserNotFound
		}
		return model.User{}, fmt.Errorf("get user by email error: %w", err)
	}

	return user, nil
}
