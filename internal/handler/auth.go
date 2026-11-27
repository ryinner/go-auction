package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/ryinner/go-auction/internal/model"
	"github.com/ryinner/go-auction/internal/store"
	"golang.org/x/crypto/bcrypt"
)

const defaultTimeout = 5 * time.Second
const passwordMinLength = 8
const passwordMaxLength = 72
const loginErrorMsg = "invalid email or password"

type UserStore interface {
	CreateUser(ctx context.Context, email string, passwordHash string) (model.User, error)
	GetUserByEmail(ctx context.Context, email string) (model.User, error)
}

type TokenGenerator interface {
	Generate(id int, email string) (string, error)
}

type Auth struct {
	responder
	store  UserStore
	tokens TokenGenerator
}

func NewAuth(s UserStore, tokens TokenGenerator, log *slog.Logger) *Auth {
	return &Auth{responder: responder{log: log}, store: s, tokens: tokens}
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

func (a *Auth) Register(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), defaultTimeout)
	defer cancel()
	var req RegisterRequest

	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		a.log.WarnContext(ctx, "register invalid json", "err", err)
		a.Error(w, r, http.StatusBadRequest, "invalid json")
		return
	}

	var err error

	passwordLen := len(req.Password)

	if passwordLen < passwordMinLength || passwordLen > passwordMaxLength {
		err = errors.Join(err, fmt.Errorf("password length must be between %d and %d", passwordMinLength, passwordMaxLength))
	}

	if !strings.Contains(req.Email, "@") || len(req.Email) > 255 {
		err = errors.Join(err, errors.New("invalid email format"))
	}

	if err != nil {
		a.Error(w, r, http.StatusBadRequest, err.Error())
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	if err != nil {
		a.Internal(w, r, "password hashing error", err)
		return
	}

	user, err := a.store.CreateUser(ctx, req.Email, string(passwordHash))

	if err != nil {
		if errors.Is(err, store.ErrEmailTaken) {
			a.Error(w, r, http.StatusConflict, "user already exists")
			return
		}
		a.Internal(w, r, "create user error", err)
		return
	}

	token, err := a.tokens.Generate(user.ID, user.Email)
	if err != nil {
		a.Internal(w, r, "generate token error", err)
		return
	}

	a.JSON(w, r, http.StatusCreated, TokenResponse{
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresIn:   int(accessTokenTTL.Seconds()),
	})
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// @todo спросить у преподавателя как правильно объединять такие похожие func
func (a *Auth) Login(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), defaultTimeout)
	defer cancel()

	var req LoginRequest

	// @todo: Вынести в request.go (parseJSON), посмотреть как работает каст типов
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		a.log.WarnContext(ctx, "login invalid json", "err", err)
		a.Error(w, r, http.StatusBadRequest, "invalid json")
		return
	}

	var err error

	if req.Password == "" {
		err = errors.Join(err, errors.New("password is required"))
	} else if len(req.Password) > passwordMaxLength {
		err = errors.Join(err, errors.New("password too long"))
	}

	if !strings.Contains(req.Email, "@") || len(req.Email) > 255 {
		err = errors.Join(err, errors.New("invalid email format"))
	}

	if err != nil {
		a.Error(w, r, http.StatusBadRequest, err.Error())
		return
	}

	user, err := a.store.GetUserByEmail(ctx, req.Email)

	if err != nil {
		if errors.Is(err, store.ErrUserNotFound) {
			a.Error(w, r, http.StatusUnauthorized, loginErrorMsg)
			return
		}
		a.Internal(w, r, "get user by email", err)
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			a.Error(w, r, http.StatusUnauthorized, loginErrorMsg)
			return
		}
		a.Internal(w, r, "compare password hash", err)
		return
	}

	token, err := a.tokens.Generate(user.ID, user.Email)
	if err != nil {
		a.Internal(w, r, "generate token error", err)
		return
	}

	a.JSON(w, r, http.StatusOK, TokenResponse{
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresIn:   int(accessTokenTTL.Seconds()),
	})
}
