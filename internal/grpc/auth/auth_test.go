package authgrpc

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/vindosVP/go-pass/internal/grpc/auth/mocks"
	"github.com/vindosVP/go-pass/internal/jwt"
	"github.com/vindosVP/go-pass/internal/models"
	authv1 "github.com/vindosVP/go-pass/internal/proto/auth"
	"github.com/vindosVP/go-pass/internal/services/auth"
	"github.com/vindosVP/go-pass/pkg/logger/sl"
)

const (
	ttl    = time.Duration(1) * time.Hour
	secret = "supersecret"
)

func TestServer_Login(t *testing.T) {

	type authMock struct {
		user   *models.User
		err    error
		needed bool
	}

	uerr := errors.New("unexpected")

	tests := []struct {
		name string
		in   *authv1.LoginRequest
		err  error
		am   authMock
	}{
		{
			name: "ok",
			in: &authv1.LoginRequest{
				Email:    "test@example.com",
				Password: "password",
			},
			err: nil,
			am: authMock{
				needed: true,
				user: &models.User{
					Email: "test@example.com",
					ID:    1,
				},
				err: nil,
			},
		},
		{
			name: "no email",
			in: &authv1.LoginRequest{
				Email:    "",
				Password: "password",
			},
			err: status.Error(codes.InvalidArgument, "email is requiered"),
			am: authMock{
				needed: false,
			},
		},
		{
			name: "wrong email",
			in: &authv1.LoginRequest{
				Email:    "wrong-email",
				Password: "password",
			},
			err: status.Error(codes.InvalidArgument, "email is invalid"),
			am: authMock{
				needed: false,
			},
		},
		{
			name: "no password",
			in: &authv1.LoginRequest{
				Email:    "test@example.com",
				Password: "",
			},
			err: status.Error(codes.InvalidArgument, "password is required"),
			am: authMock{
				needed: false,
			},
		},
		{
			name: "invalid email or password",
			in: &authv1.LoginRequest{
				Email:    "test@example.com",
				Password: "password",
			},
			err: status.Error(codes.Unauthenticated, "invalid email or password"),
			am: authMock{
				needed: true,
				user:   nil,
				err:    auth.ErrInvalidCredentials,
			},
		},
		{
			name: "unexpected error",
			in: &authv1.LoginRequest{
				Email:    "test@example.com",
				Password: "password",
			},
			err: status.Error(codes.Internal, "failed to login"),
			am: authMock{
				needed: true,
				user:   nil,
				err:    uerr,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sl.SetupLogger("test")
			a := mocks.NewAuth(t)
			if tt.am.needed {
				token := ""
				if tt.am.err != nil {
					token = ""
				} else {
					tok, err := jwt.NewToken(tt.am.user, ttl, secret)
					require.NoError(t, err)
					token = tok
				}
				a.On("Login", mock.Anything, tt.in.Email, tt.in.Password).Return(token, tt.am.err)
			}
			s := server{
				auth: a,
			}
			out, err := s.Login(context.Background(), tt.in)
			assert.ErrorIs(t, err, tt.err)
			if err == nil {
				require.NotEqual(t, "", out.Token)
				email, err := jwt.VerifyToken(out.Token, secret)
				require.NoError(t, err)
				assert.Equal(t, tt.in.Email, email)
			}
		})
	}
}

func TestServer_Register(t *testing.T) {

	type authMock struct {
		user   *models.User
		err    error
		needed bool
	}

	uerr := errors.New("unexpected")

	tests := []struct {
		name string
		in   *authv1.RegisterRequest
		w    *authv1.RegisterResponse
		err  error
		am   authMock
	}{
		{
			name: "ok",
			in: &authv1.RegisterRequest{
				Email:    "test@example.com",
				Password: "password",
			},
			w: &authv1.RegisterResponse{
				UserId: 1,
			},
			err: nil,
			am: authMock{
				needed: true,
				user: &models.User{
					Email: "test@example.com",
					ID:    1,
				},
				err: nil,
			},
		},
		{
			name: "no email",
			in: &authv1.RegisterRequest{
				Email:    "",
				Password: "password",
			},
			w:   nil,
			err: status.Error(codes.InvalidArgument, "email is requiered"),
			am: authMock{
				needed: false,
			},
		},
		{
			name: "wrong email",
			in: &authv1.RegisterRequest{
				Email:    "wrong-email",
				Password: "password",
			},
			w:   nil,
			err: status.Error(codes.InvalidArgument, "email is invalid"),
			am: authMock{
				needed: false,
			},
		},
		{
			name: "no password",
			in: &authv1.RegisterRequest{
				Email:    "test@example.com",
				Password: "",
			},
			w:   nil,
			err: status.Error(codes.InvalidArgument, "password is required"),
			am: authMock{
				needed: false,
			},
		},
		{
			name: "user already exists",
			in: &authv1.RegisterRequest{
				Email:    "test@example.com",
				Password: "password",
			},
			w:   nil,
			err: status.Error(codes.InvalidArgument, "user with this email already exists"),
			am: authMock{
				needed: true,
				user:   nil,
				err:    auth.ErrUserAlreadyExists,
			},
		},
		{
			name: "unexpected error",
			in: &authv1.RegisterRequest{
				Email:    "test@example.com",
				Password: "password",
			},
			w:   nil,
			err: status.Error(codes.Internal, "failed to create user"),
			am: authMock{
				needed: true,
				user:   nil,
				err:    uerr,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sl.SetupLogger("test")
			a := mocks.NewAuth(t)
			if tt.am.needed {
				a.On("CreateUser", mock.Anything, tt.in.Email, tt.in.Password).Return(tt.am.user, tt.am.err)
			}
			s := server{
				auth: a,
			}
			out, err := s.Register(context.Background(), tt.in)
			assert.ErrorIs(t, err, tt.err)
			if err == nil {
				assert.Equal(t, tt.w, out)
			}
		})
	}
}
