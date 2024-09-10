// Package auth provides an API to create and log in users.
package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/vindosVP/go-pass/internal/jwt"
	"github.com/vindosVP/go-pass/internal/models"
	"github.com/vindosVP/go-pass/internal/storage"
	"github.com/vindosVP/go-pass/pkg/logger/sl"
)

// UserStorage is a user storage interface.
//
//go:generate go run github.com/vektra/mockery/v2@v2.43.2 --name=UserStorage
type UserStorage interface {
	CreateUser(ctx context.Context, email string, passHash []byte) (*models.User, error)
	UserByEmail(ctx context.Context, email string) (*models.User, error)
}

var (
	// ErrUserAlreadyExists - user already exists error.
	ErrUserAlreadyExists = errors.New("user already exists")

	// ErrInvalidCredentials - invalid credentials error.
	ErrInvalidCredentials = errors.New("invalid credentials")
)

// Auth consists the authentication fields.
type Auth struct {
	userStorage UserStorage
	tokenTTL    time.Duration
	secret      string
}

// New creates the Auth instance.
func New(us UserStorage, secret string) *Auth {
	return &Auth{userStorage: us, secret: secret}
}

// CreateUser creates a new user with provided email and password.
func (a *Auth) CreateUser(ctx context.Context, email string, pass string) (*models.User, error) {

	lg := sl.Log.With(slog.String("email", email))
	lg.Info("creating user")

	passHash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		lg.Error("failed to generate password hash", sl.Err(err))
		return nil, fmt.Errorf("failed to generate password hash: %w", err)
	}
	user, err := a.userStorage.CreateUser(ctx, email, passHash)
	if err != nil {
		if errors.Is(err, storage.ErrUserAlreadyExists) {
			lg.Info("user already exists")
			return nil, ErrUserAlreadyExists
		}
		lg.Error("failed to create user", sl.Err(err))
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	lg.Info("user created")
	return user, nil
}

// Login logs in user with provided email and password.
func (a *Auth) Login(ctx context.Context, email string, pass string) (string, error) {

	lg := sl.Log.With(slog.String("email", email))
	lg.Info("logging in user")

	user, err := a.userStorage.UserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotExist) {
			lg.Info("invalid credentials")
			return "", ErrInvalidCredentials
		}
		lg.Error("failed to get user", sl.Err(err))
		return "", fmt.Errorf("failed to get user: %w", err)
	}
	err = bcrypt.CompareHashAndPassword(user.PassHash, []byte(pass))
	if err != nil {
		lg.Info("invalid credentials")
		return "", ErrInvalidCredentials
	}
	token, err := jwt.NewToken(user, a.secret)
	if err != nil {
		lg.Error("failed to create token", sl.Err(err))
		return "", fmt.Errorf("failed to create token: %w", err)
	}

	lg.Info("user logged in")
	return token, nil
}
