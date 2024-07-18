package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/jackc/pgx/v5/pgxpool"

	serverConfig "github.com/vindosVP/go-pass/cmd/server/config"
	"github.com/vindosVP/go-pass/internal/app"
	"github.com/vindosVP/go-pass/pkg/db"
	"github.com/vindosVP/go-pass/pkg/logger/sl"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func main() {
	ctx := context.Background()
	printBuildInfo()
	conf := serverConfig.MustLoad()
	sl.SetupLogger(conf.Env)

	sl.Log.Info("Starting server...", slog.String("config", conf.String()))

	dsn := db.PostgresDSN(conf.DB.Host, conf.DB.Port, conf.DB.User, conf.DB.Password, conf.DB.Database)
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatal(fmt.Errorf("error connecting to database: %w", err))
	}
	a := app.New(conf.GRPC.Port, pool, conf.Auth.TokenTTL, conf.Auth.Secret)

	go func() {
		a.MustRun()
	}()

	stop := make(chan os.Signal, 3)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	<-stop

	a.Stop()
	sl.Log.Info("Gracefully stopped")
}

func printBuildInfo() {
	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)
}
