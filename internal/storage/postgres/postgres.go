// Package postgres is a package for postgres storage
package postgres

import (
	"context"
	"errors"
	"fmt"
	"syscall"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/vindosVP/go-pass/internal/models"
	"github.com/vindosVP/go-pass/internal/storage"
	"github.com/vindosVP/go-pass/pkg/logger/sl"
)

var retryDelays = map[uint]time.Duration{
	0: 1 * time.Second,
	1: 3 * time.Second,
	2: 5 * time.Second,
}

// Storage consists the database
type Storage struct {
	db *pgxpool.Pool
}

// New creates the Storage instance
func New(db *pgxpool.Pool) *Storage {
	return &Storage{db: db}
}

// CreateUser creates user with provided email and password hash
func (s *Storage) CreateUser(ctx context.Context, email string, passHash []byte) (*models.User, error) {
	return retry.DoWithData(func() (*models.User, error) {
		query := "insert into users (email, hashed_password, created_at) values ($1, $2, $3) returning id, email, hashed_password, created_at"
		row := s.db.QueryRow(ctx, query, email, passHash, time.Now())
		user := &models.User{}
		if err := row.Scan(&user.ID, &user.Email, &user.PassHash, &user.CreatedAt); err != nil {
			if pgErrCode(err) == pgerrcode.UniqueViolation {
				return nil, storage.ErrUserAlreadyExists
			}
			return nil, err
		}
		return user, nil
	}, retryOpts()...)
}

// UserByEmail finds a user by provided email
func (s *Storage) UserByEmail(ctx context.Context, email string) (*models.User, error) {
	return retry.DoWithData(func() (*models.User, error) {
		query := "select id, email, hashed_password, created_at from users where email = $1"
		row := s.db.QueryRow(ctx, query, email)
		user := &models.User{}
		if err := row.Scan(&user.ID, &user.Email, &user.PassHash, &user.CreatedAt); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, storage.ErrUserNotExist
			}
			return nil, err
		}
		return user, nil
	}, retryOpts()...)
}

func retryOpts() []retry.Option {
	return []retry.Option{
		retry.RetryIf(func(err error) bool {
			return pgerrcode.IsConnectionException(pgErrCode(err)) || errors.Is(err, syscall.ECONNREFUSED)
		}),
		retry.DelayType(func(n uint, err error, config *retry.Config) time.Duration {
			delay := retryDelays[n]
			return delay
		}),
		retry.OnRetry(func(n uint, err error) {
			sl.Log.Info(fmt.Sprintf("Failed to connect to database, retrying in %s", retryDelays[n]))
		}),
		retry.Attempts(4),
		retry.LastErrorOnly(true),
	}
}

func pgErrCode(err error) string {
	if e, ok := err.(*pgconn.PgError); ok {
		return e.Code
	}

	return ""
}
