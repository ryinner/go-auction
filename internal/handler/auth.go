package handler

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/ryinner/go-auction/internal/model"
	"github.com/ryinner/go-auction/internal/store"
	"golang.org/x/crypto/bcrypt"
)

const defaultTimeout = 5 * time.Second

type UserStore interface {
	CreateUser(ctx context.Context, email string, passwordHash string) (model.User, error)
	GetUserByEmail(ctx context.Context, email string) (model.User, error)
}

type Auth struct {
	jwtSecret string
	store     UserStore
	log       *slog.Logger
}

func newAuth(s UserStore, log *slog.Logger, jwtSecret string) *Auth {
	return &Auth{store: s, log: log, jwtSecret: jwtSecret}
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (a *Auth) Register(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), defaultTimeout)
	defer cancel()
	var req RegisterRequest

	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		a.log.Warn("register invalid json", "err", err)
		err = jsonResponse(w, map[string]string{
			"error": "Invalid json",
		}, http.StatusInternalServerError)
		return
	}

	var err error

	passwordLen := len(req.Password)

	if passwordLen < 8 || passwordLen > 72 {
		err = errors.Join(err, errors.New("password length must be between 8 and 72"))
	}

	if !strings.Contains(req.Email, "@") || len(req.Email) > 255 {
		err = errors.Join(err, errors.New("invalid email format"))
	}

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		err = jsonResponse(w, map[string]string{
			"error": err.Error(),
		}, http.StatusBadRequest)
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	if err != nil {
		a.log.ErrorContext(ctx, "password hashing error", "err", err)
		err = jsonResponse(w, map[string]string{
			"error": "Interal server error",
		}, http.StatusInternalServerError)
		return
	}

	user, err := a.store.CreateUser(ctx, req.Email, string(passwordHash))

	if err != nil {
		if errors.Is(err, store.ErrEmailTaken) {
			err = jsonResponse(w, map[string]string{
				"error": "User already exists",
			}, http.StatusConflict)
			return
		}
		a.log.ErrorContext(ctx, "create user error", "err", err)
		err = jsonResponse(w, map[string]string{
			"error": "Internal server error",
		}, http.StatusInternalServerError)
		return
	}

	// jwt token service
}
