// Package authgrpc consists the auth grpc server
package authgrpc

import (
	"context"
	"errors"
	"log/slog"
	"net/mail"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/vindosVP/go-pass/internal/models"
	authv1 "github.com/vindosVP/go-pass/internal/proto/auth"
	"github.com/vindosVP/go-pass/internal/services/auth"
	"github.com/vindosVP/go-pass/pkg/logger/sl"
)

// Auth is an authentication API interface
//
//go:generate go run github.com/vektra/mockery/v2@v2.43.2 --name=Auth
type Auth interface {
	CreateUser(ctx context.Context, email string, pass string) (*models.User, error)
	Login(ctx context.Context, email string, pass string) (string, error)
}

type server struct {
	authv1.UnimplementedAuthServer
	auth Auth
}

// Register registers the auth service.
func Register(gRPCServer *grpc.Server, auth Auth) {
	authv1.RegisterAuthServer(gRPCServer, &server{auth: auth})
}

// Login logs in the user
func (s *server) Login(ctx context.Context, in *authv1.LoginRequest) (*authv1.LoginResponse, error) {

	lg := sl.Log.With(slog.String("email", in.Email))
	lg.Info("handling login")

	err, code, msg := validate(in.Email, in.Password)
	if err != nil {
		lg.Info(msg)
		return nil, status.Error(code, msg)
	}
	token, err := s.auth.Login(ctx, in.Email, in.Password)
	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			lg.Info("invalid email or password")
			return nil, status.Error(codes.Unauthenticated, "invalid email or password")
		}
		lg.Error("failed to login", sl.Err(err))
		return nil, status.Error(codes.Internal, "failed to login")
	}

	lg.Info("logged in")
	return &authv1.LoginResponse{Token: token}, nil
}

// Register registers a new user
func (s *server) Register(ctx context.Context, in *authv1.RegisterRequest) (*authv1.RegisterResponse, error) {

	lg := sl.Log.With(slog.String("email", in.Email))
	lg.Info("handling register")

	err, code, msg := validate(in.Email, in.Password)
	if err != nil {
		lg.Info(msg)
		return nil, status.Error(code, msg)
	}
	user, err := s.auth.CreateUser(ctx, in.Email, in.Password)
	if err != nil {
		if errors.Is(err, auth.ErrUserAlreadyExists) {
			lg.Info("user with this email already exists")
			return nil, status.Error(codes.InvalidArgument, "user with this email already exists")
		}
		lg.Error("failed to create user", sl.Err(err))
		return nil, status.Error(codes.Internal, "failed to create user")
	}

	lg.Info("registered user")
	return &authv1.RegisterResponse{UserId: int64(user.ID)}, nil
}

func validate(email string, password string) (error, codes.Code, string) {
	errValidation := errors.New("user validation error")
	if email == "" {
		return errValidation, codes.InvalidArgument, "email is requiered"
	}
	if !isValidEmail(email) {
		return errValidation, codes.InvalidArgument, "email is invalid"
	}
	if password == "" {
		return errValidation, codes.InvalidArgument, "password is required"
	}
	return nil, codes.OK, ""
}

func isValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
