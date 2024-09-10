package main

import (
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/golang-migrate/migrate/v4/source/github"
	_ "github.com/mattes/migrate/source/file"

	migratorConfig "github.com/vindosVP/go-pass/cmd/migrator/config"
	"github.com/vindosVP/go-pass/pkg/db"
)

func main() {

	conf := migratorConfig.MustLoad()
	m, err := migrate.New(
		fmt.Sprintf("file://%s", conf.MigrationsPath),
		db.PostgresDSN(conf.DB.Host, conf.DB.Port, conf.DB.User, conf.DB.Password, conf.DB.Database))

	if err != nil {
		panic(fmt.Errorf("failed to connect database: %w", err))
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("no migrations to apply")
			return
		}
		panic(fmt.Errorf("failed to apply migrations: %w", err))
	}

	fmt.Println("migrations applied")
}
