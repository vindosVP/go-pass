// Package app creates and works with the App
package app

import (
	"github.com/jackc/pgx/v5/pgxpool"

	grpcapp "github.com/vindosVP/go-pass/internal/app/grpc"
	"github.com/vindosVP/go-pass/internal/services/auth"
	"github.com/vindosVP/go-pass/internal/services/passkeeper"
	"github.com/vindosVP/go-pass/internal/storage/postgres"
)

// App consist the grpc server
type App struct {
	grpcServer *grpcapp.App
}

// MustRun runs the app
func (a *App) MustRun() {
	a.grpcServer.MustRun()
}

// Stop stops app
func (a *App) Stop() {
	a.grpcServer.Stop()
}

// New creates the App instance
func New(port int, pool *pgxpool.Pool, secret string, fl string) *App {
	s := postgres.New(pool)
	a := auth.New(s, secret)
	k := passkeeper.New(s, s, s, s, fl)
	grpcApp := grpcapp.New(port, secret, a, k)
	return &App{
		grpcServer: grpcApp,
	}
}
