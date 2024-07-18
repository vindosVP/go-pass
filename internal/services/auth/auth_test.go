package auth

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"

	"github.com/vindosVP/go-pass/internal/jwt"
	"github.com/vindosVP/go-pass/internal/models"
	"github.com/vindosVP/go-pass/internal/services/auth/mocks"
	"github.com/vindosVP/go-pass/internal/storage"
	"github.com/vindosVP/go-pass/pkg/logger/sl"
)

const (
	ttl    = time.Duration(1) * time.Hour
	secret = "supersecret"
)

func newUser(email string, pass string) *models.User {
	passHash, _ := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	return &models.User{
		ID:        1,
		Email:     email,
		PassHash:  passHash,
		CreatedAt: time.Time{},
	}
}

type storageMock struct {
	user *models.User
	err  error
}

func TestAuth_Login(t *testing.T) {

	type want struct {
		checkToken bool
		err        error
	}

	type fields struct {
		email    string
		password string
	}

	type storageMock struct {
		user *models.User
		err  error
	}

	tests := []struct {
		name string
		sm   storageMock
		f    fields
		w    want
	}{
		{
			name: "ok",
			sm: storageMock{
				user: newUser("test@test.com", "password"),
				err:  nil,
			},
			f: fields{
				email:    "test@test.com",
				password: "password",
			},
			w: want{
				checkToken: true,
				err:        nil,
			},
		},
		{
			name: "user not found",
			sm: storageMock{
				user: nil,
				err:  storage.ErrUserNotExist,
			},
			f: fields{
				email:    "test@test.com",
				password: "password",
			},
			w: want{
				checkToken: false,
				err:        ErrInvalidCredentials,
			},
		},
		{
			name: "invalid password",
			sm: storageMock{
				user: newUser("test@test.com", "password"),
				err:  nil,
			},
			f: fields{
				email:    "test@test.com",
				password: "wrong-password",
			},
			w: want{
				checkToken: false,
				err:        ErrInvalidCredentials,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sl.SetupLogger("test")
			ms := mocks.NewUserStorage(t)
			ms.On("UserByEmail", mock.Anything, tt.f.email).Return(tt.sm.user, tt.sm.err)
			a := New(ms, ttl, secret)
			token, err := a.Login(context.Background(), tt.f.email, tt.f.password)
			require.ErrorIs(t, err, tt.w.err)
			if tt.w.checkToken {
				require.NotEqual(t, "", token)
				email, err := jwt.VerifyToken(token, secret)
				require.NoError(t, err)
				assert.Equal(t, tt.f.email, email)
			}
		})
	}

}

func TestAuth_CreateUser(t *testing.T) {

	type want struct {
		user *models.User
		err  error
	}

	type fields struct {
		email string
		pass  string
	}

	unexpectedErr := errors.New("unexpected")

	tests := []struct {
		name string
		sm   storageMock
		f    fields
		w    want
	}{
		{
			name: "ok",
			sm: storageMock{
				user: newUser("test@example.com", "password"),
				err:  nil,
			},
			f: fields{
				email: "test@example.com",
				pass:  "password",
			},
			w: want{
				user: newUser("test@example.com", "password"),
				err:  nil,
			},
		},
		{
			name: "user already exists",
			sm: storageMock{
				user: nil,
				err:  storage.ErrUserAlreadyExists,
			},
			f: fields{
				email: "test@example.com",
				pass:  "password",
			},
			w: want{
				user: nil,
				err:  ErrUserAlreadyExists,
			},
		},
		{
			name: "unexpected error",
			sm: storageMock{
				user: nil,
				err:  unexpectedErr,
			},
			f: fields{
				email: "test@example.com",
				pass:  "password",
			},
			w: want{
				user: nil,
				err:  unexpectedErr,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sl.SetupLogger("test")
			ms := mocks.NewUserStorage(t)
			ms.On("CreateUser", mock.Anything, tt.f.email, mock.Anything).Return(tt.sm.user, tt.sm.err)
			a := New(ms, ttl, secret)
			usr, err := a.CreateUser(context.Background(), tt.f.email, tt.f.pass)
			if tt.w.user != nil {
				assert.Equal(t, tt.w.user.Email, usr.Email)
				assert.NoError(t, bcrypt.CompareHashAndPassword(usr.PassHash, []byte(tt.f.pass)))
			}
			assert.ErrorIs(t, err, tt.w.err)
		})
	}

}
