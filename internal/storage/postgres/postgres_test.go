package postgres

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/golang-migrate/migrate/v4"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/vindosVP/go-pass/internal/models"
	"github.com/vindosVP/go-pass/internal/storage"
	"github.com/vindosVP/go-pass/pkg/db"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type testConfig struct {
	Host           string `yaml:"host" validate:"required"`
	Port           int    `yaml:"port" validate:"required"`
	User           string `yaml:"user" validate:"required"`
	Password       string `yaml:"password" validate:"required"`
	Database       string `yaml:"database" validate:"required"`
	MigrationsPath string `yaml:"migrationsPath" validate:"required"`
}

func parseTestConfig() (*testConfig, error) {
	testCfgPath := "../../../configs/test-config.yaml"
	viper.SetConfigFile(testCfgPath)
	conf := &testConfig{}

	err := viper.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	err = viper.Unmarshal(conf)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	validate := validator.New()
	if err := validate.Struct(conf); err != nil {
		panic(fmt.Errorf("missing requiered attributes: %w", err))
	}

	return conf, nil
}

func cleanup(cfg *testConfig) {

	m, _ := migrate.New(
		fmt.Sprintf("file://%s", cfg.MigrationsPath),
		db.PostgresDSN(cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Database))

	err := m.Down()
	if err != nil {
		panic(err)
	}

}

func setupDatabase(cfg *testConfig) error {

	m, err := migrate.New(
		fmt.Sprintf("file://%s", cfg.MigrationsPath),
		db.PostgresDSN(cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Database))

	if err != nil {
		panic(fmt.Errorf("failed to connect database: %w", err))
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			return nil
		}
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	return nil

}

func TestStorage_CreateUser(t *testing.T) {

	ctx := context.Background()

	cfg, err := parseTestConfig()
	require.NoError(t, err)

	err = setupDatabase(cfg)
	require.NoError(t, err)
	defer cleanup(cfg)

	pool, err := pgxpool.New(ctx, db.PostgresDSN(cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Database))
	require.NoError(t, err)

	s := New(pool)
	user := &models.User{
		ID:       1,
		Email:    "testmail@gmail.com",
		PassHash: []byte("supersecretpassword"),
	}

	res, err := s.CreateUser(ctx, user.Email, user.PassHash)
	assert.NoError(t, err)
	assert.Equal(t, user.ID, res.ID)
	assert.Equal(t, user.Email, res.Email)
	assert.Equal(t, user.PassHash, res.PassHash)

	res, err = s.CreateUser(ctx, user.Email, user.PassHash)
	assert.ErrorIs(t, err, storage.ErrUserAlreadyExists)

}

func TestStorage_UserByEmail(t *testing.T) {

	ctx := context.Background()

	cfg, err := parseTestConfig()
	require.NoError(t, err)

	err = setupDatabase(cfg)
	require.NoError(t, err)
	defer cleanup(cfg)

	pool, err := pgxpool.New(ctx, db.PostgresDSN(cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Database))
	require.NoError(t, err)

	s := New(pool)
	user := &models.User{
		ID:       1,
		Email:    "testmail@gmail.com",
		PassHash: []byte("supersecretpassword"),
	}

	_, err = s.UserByEmail(ctx, user.Email)
	assert.ErrorIs(t, err, storage.ErrUserNotExist)

	_, err = s.CreateUser(ctx, user.Email, user.PassHash)
	require.NoError(t, err)

	res, err := s.UserByEmail(ctx, user.Email)
	assert.NoError(t, err)
	assert.Equal(t, user.ID, res.ID)
	assert.Equal(t, user.Email, res.Email)
	assert.Equal(t, user.PassHash, res.PassHash)

}
