package main

import (
	"fmt"
	"log/slog"

	serverConfig "github.com/vindosVP/go-pass/cmd/server/config"
	"github.com/vindosVP/go-pass/pkg/logger/sl"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func main() {
	printBuildInfo()
	conf := serverConfig.MustLoad()
	sl.SetupLogger(conf.Env)

	sl.Log.Info("Starting server...", slog.String("config", conf.String()))

}

func printBuildInfo() {
	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)
}
