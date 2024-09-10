package postgres

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

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
	testCfgPath := "C:/Golang/github.com/vindosVp/go-pass/configs/test-config.yaml"
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

func setupStorage(cfg *testConfig, createUser bool) (*Storage, error) {

	ctx := context.Background()

	m, err := migrate.New(
		fmt.Sprintf("file://%s", cfg.MigrationsPath),
		db.PostgresDSN(cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Database))

	if err != nil {
		panic(fmt.Errorf("failed to connect database: %w", err))
	}

	if err := m.Up(); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			return nil, fmt.Errorf("failed to apply migrations: %w", err)
		}
	}

	pool, err := pgxpool.New(ctx, db.PostgresDSN(cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Database))
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	s := New(pool)
	if createUser {
		user := &models.User{
			ID:       1,
			Email:    "testmail@gmail.com",
			PassHash: []byte("supersecretpassword"),
		}

		_, err = s.CreateUser(ctx, user.Email, user.PassHash)
		if err != nil {
			return nil, fmt.Errorf("failed to create user: %w", err)
		}
	}

	return s, nil

}

func TestStorage_AddPassword(t *testing.T) {
	ctx := context.Background()

	cfg, err := parseTestConfig()
	require.NoError(t, err)

	s, err := setupStorage(cfg, true)
	defer cleanup(cfg)
	require.NoError(t, err)

	pwd := &models.Password{
		ID:       1,
		OwnerID:  1,
		Login:    "login",
		Password: "password",
		Metadata: "metadata",
	}
	id, err := s.AddPassword(ctx, pwd)
	require.NoError(t, err)
	assert.Equal(t, pwd.ID, id)
}

func TestStorage_AddCard(t *testing.T) {
	ctx := context.Background()

	cfg, err := parseTestConfig()
	require.NoError(t, err)

	s, err := setupStorage(cfg, true)
	defer cleanup(cfg)
	require.NoError(t, err)

	card := &models.Card{
		ID:       1,
		OwnerID:  1,
		Number:   "1234 1234 1234 1234",
		CVC:      "123",
		Owner:    "USER CARD",
		Date:     "06/28",
		Metadata: "metadata",
	}
	id, err := s.AddCard(ctx, card)
	require.NoError(t, err)
	assert.Equal(t, card.ID, id)
}

func TestStorage_AddText(t *testing.T) {
	ctx := context.Background()

	cfg, err := parseTestConfig()
	require.NoError(t, err)

	s, err := setupStorage(cfg, true)
	defer cleanup(cfg)
	require.NoError(t, err)

	text := &models.Text{
		ID:       1,
		OwnerID:  1,
		Text:     "some text",
		Metadata: "metadata",
	}
	id, err := s.AddText(ctx, text)
	require.NoError(t, err)
	assert.Equal(t, text.ID, id)
}

func TestStorage_AddFile(t *testing.T) {
	ctx := context.Background()

	cfg, err := parseTestConfig()
	require.NoError(t, err)

	s, err := setupStorage(cfg, true)
	defer cleanup(cfg)
	require.NoError(t, err)

	file := &models.File{
		ID:       1,
		OwnerID:  1,
		FileName: "file.txt",
		Metadata: "metadata",
	}
	id, err := s.AddFile(ctx, file)
	require.NoError(t, err)
	assert.Equal(t, file.ID, id)
}

func TestStorage_UpdatePassword(t *testing.T) {

	ctx := context.Background()

	cfg, err := parseTestConfig()
	require.NoError(t, err)

	s, err := setupStorage(cfg, true)
	defer cleanup(cfg)
	require.NoError(t, err)

	pwd := &models.Password{
		ID:       1,
		OwnerID:  1,
		Login:    "login",
		Password: "password",
		Metadata: "metadata",
	}
	_, err = s.AddPassword(ctx, pwd)
	require.NoError(t, err)

	pwd.Login = "new-login"
	pwd.Password = "new-password"
	pwd.Metadata = "new-metadata"

	err = s.UpdatePassword(ctx, pwd)
	assert.NoError(t, err)
}

func TestStorage_UpdateCard(t *testing.T) {

	ctx := context.Background()

	cfg, err := parseTestConfig()
	require.NoError(t, err)

	s, err := setupStorage(cfg, true)
	defer cleanup(cfg)
	require.NoError(t, err)

	card := &models.Card{
		ID:       1,
		OwnerID:  1,
		Number:   "1234 1234 1234 1234",
		CVC:      "123",
		Owner:    "USER CARD",
		Date:     "06/28",
		Metadata: "metadata",
	}
	_, err = s.AddCard(ctx, card)
	require.NoError(t, err)

	card.Number = "2222 2222 2222 2222"
	card.CVC = "321"
	card.Owner = "CARD OWNER"
	card.Date = "01/20"
	card.Metadata = "new-metadata"

	err = s.UpdateCard(ctx, card)
	assert.NoError(t, err)
}

func TestStorage_UpdateText(t *testing.T) {
	ctx := context.Background()

	cfg, err := parseTestConfig()
	require.NoError(t, err)

	s, err := setupStorage(cfg, true)
	defer cleanup(cfg)
	require.NoError(t, err)

	text := &models.Text{
		ID:       1,
		OwnerID:  1,
		Text:     "some text",
		Metadata: "metadata",
	}
	_, err = s.AddText(ctx, text)
	require.NoError(t, err)

	text.Text = "some updated text"
	text.Metadata = "new metadata"

	err = s.UpdateText(ctx, text)
	assert.NoError(t, err)
}

func TestStorage_DeletePassword(t *testing.T) {

	ctx := context.Background()

	cfg, err := parseTestConfig()
	require.NoError(t, err)

	s, err := setupStorage(cfg, true)
	defer cleanup(cfg)
	require.NoError(t, err)

	pwd := &models.Password{
		ID:       1,
		OwnerID:  1,
		Login:    "login",
		Password: "password",
		Metadata: "metadata",
	}
	id, err := s.AddPassword(ctx, pwd)
	require.NoError(t, err)

	err = s.DeletePassword(ctx, id, pwd.OwnerID)
	assert.NoError(t, err)
}

func TestStorage_DeleteCard(t *testing.T) {

	ctx := context.Background()

	cfg, err := parseTestConfig()
	require.NoError(t, err)

	s, err := setupStorage(cfg, true)
	defer cleanup(cfg)
	require.NoError(t, err)

	card := &models.Card{
		ID:       1,
		OwnerID:  1,
		Number:   "1234 1234 1234 1234",
		CVC:      "123",
		Owner:    "USER CARD",
		Date:     "06/28",
		Metadata: "metadata",
	}
	id, err := s.AddCard(ctx, card)
	require.NoError(t, err)

	err = s.DeleteCard(ctx, id, card.OwnerID)
	assert.NoError(t, err)
}

func TestStorage_DeleteText(t *testing.T) {
	ctx := context.Background()

	cfg, err := parseTestConfig()
	require.NoError(t, err)

	s, err := setupStorage(cfg, true)
	defer cleanup(cfg)
	require.NoError(t, err)

	text := &models.Text{
		ID:       1,
		OwnerID:  1,
		Text:     "some text",
		Metadata: "metadata",
	}
	id, err := s.AddText(ctx, text)
	require.NoError(t, err)

	err = s.DeleteText(ctx, id, text.OwnerID)
	assert.NoError(t, err)
}

func TestStorage_MarkFileAsUploaded(t *testing.T) {

	ctx := context.Background()

	cfg, err := parseTestConfig()
	require.NoError(t, err)

	s, err := setupStorage(cfg, true)
	defer cleanup(cfg)
	require.NoError(t, err)

	file := &models.File{
		ID:       1,
		OwnerID:  1,
		FileName: "file.txt",
		Metadata: "metadata",
	}
	id, err := s.AddFile(ctx, file)
	require.NoError(t, err)

	err = s.MarkFileAsUploaded(ctx, id, file.OwnerID)
	assert.NoError(t, err)
}

func TestStorage_GetFile(t *testing.T) {

	ctx := context.Background()

	cfg, err := parseTestConfig()
	require.NoError(t, err)

	s, err := setupStorage(cfg, true)
	defer cleanup(cfg)
	require.NoError(t, err)

	file := &models.File{
		ID:       1,
		OwnerID:  1,
		FileName: "file.txt",
		Metadata: "metadata",
	}

	_, err = s.GetFile(ctx, file.ID, file.OwnerID)
	assert.ErrorIs(t, err, storage.ErrFileNotExist)

	id, err := s.AddFile(ctx, file)
	require.NoError(t, err)

	f, err := s.GetFile(ctx, id, file.OwnerID)
	assert.NoError(t, err)
	assert.Equal(t, file.FileName, f.FileName)
	assert.Equal(t, file.Metadata, f.Metadata)
}

func TestStorage_DeleteFile(t *testing.T) {

	ctx := context.Background()

	cfg, err := parseTestConfig()
	require.NoError(t, err)

	s, err := setupStorage(cfg, true)
	defer cleanup(cfg)
	require.NoError(t, err)

	file := &models.File{
		ID:       1,
		OwnerID:  1,
		FileName: "file.txt",
		Metadata: "metadata",
	}
	id, err := s.AddFile(ctx, file)
	require.NoError(t, err)

	err = s.DeleteFile(ctx, id, file.OwnerID)
	assert.NoError(t, err)
}

func TestStorage_GetPasswords(t *testing.T) {

	ctx := context.Background()

	cfg, err := parseTestConfig()
	require.NoError(t, err)

	s, err := setupStorage(cfg, true)
	defer cleanup(cfg)
	require.NoError(t, err)

	pwds := []*models.Password{
		{
			ID:        1,
			OwnerID:   1,
			Login:     "login1",
			Password:  "password1",
			Metadata:  "metadata1",
			CreatedAt: time.Time{},
		},
		{
			ID:        2,
			OwnerID:   1,
			Login:     "login2",
			Password:  "password2",
			Metadata:  "metadata2",
			CreatedAt: time.Time{},
		},
		{
			ID:        3,
			OwnerID:   1,
			Login:     "login3",
			Password:  "password3",
			Metadata:  "metadata3",
			CreatedAt: time.Time{},
		},
	}

	for _, pwd := range pwds {
		_, err = s.AddPassword(ctx, pwd)
		require.NoError(t, err)
	}

	res, err := s.GetPasswords(ctx, 1)
	for _, r := range res {
		r.CreatedAt = time.Time{}
	}
	assert.NoError(t, err)
	assert.ElementsMatch(t, pwds, res)
}

func TestStorage_GetCards(t *testing.T) {

	ctx := context.Background()

	cfg, err := parseTestConfig()
	require.NoError(t, err)

	s, err := setupStorage(cfg, true)
	defer cleanup(cfg)
	require.NoError(t, err)

	cards := []*models.Card{
		{
			ID:        1,
			OwnerID:   1,
			Number:    "1234",
			CVC:       "123",
			Owner:     "VADIM VALOV",
			Date:      "06/28",
			Metadata:  "md",
			CreatedAt: time.Time{},
		},
		{
			ID:        2,
			OwnerID:   1,
			Number:    "5678",
			CVC:       "111",
			Owner:     "IVANOV IVAN",
			Date:      "06/20",
			Metadata:  "md1",
			CreatedAt: time.Time{},
		},
		{
			ID:        3,
			OwnerID:   1,
			Number:    "9012",
			CVC:       "222",
			Owner:     "VASYA IVALOV",
			Date:      "01/21",
			Metadata:  "md2",
			CreatedAt: time.Time{},
		},
	}

	for _, card := range cards {
		_, err = s.AddCard(ctx, card)
		require.NoError(t, err)
	}

	res, err := s.GetCards(ctx, 1)
	for _, r := range res {
		r.CreatedAt = time.Time{}
	}
	assert.NoError(t, err)
	assert.ElementsMatch(t, cards, res)
}

func TestStorage_GetTexts(t *testing.T) {

	ctx := context.Background()

	cfg, err := parseTestConfig()
	require.NoError(t, err)

	s, err := setupStorage(cfg, true)
	defer cleanup(cfg)
	require.NoError(t, err)

	texts := []*models.Text{
		{
			ID:        1,
			OwnerID:   1,
			Text:      "text1",
			Metadata:  "md1",
			CreatedAt: time.Time{},
		},
		{
			ID:        2,
			OwnerID:   1,
			Text:      "text2",
			Metadata:  "md2",
			CreatedAt: time.Time{},
		},
		{
			ID:        3,
			OwnerID:   1,
			Text:      "text3",
			Metadata:  "md3",
			CreatedAt: time.Time{},
		},
	}

	for _, text := range texts {
		_, err = s.AddText(ctx, text)
		require.NoError(t, err)
	}

	res, err := s.GetTexts(ctx, 1)
	for _, r := range res {
		r.CreatedAt = time.Time{}
	}
	assert.NoError(t, err)
	assert.ElementsMatch(t, texts, res)
}

func TestStorage_GetFiles(t *testing.T) {

	ctx := context.Background()

	cfg, err := parseTestConfig()
	require.NoError(t, err)

	s, err := setupStorage(cfg, true)
	defer cleanup(cfg)
	require.NoError(t, err)

	files := []*models.File{
		{
			ID:        1,
			OwnerID:   1,
			FileName:  "file1",
			Metadata:  "md1",
			CreatedAt: time.Time{},
		},
		{
			ID:        2,
			OwnerID:   1,
			FileName:  "file2",
			Metadata:  "md2",
			CreatedAt: time.Time{},
		},
		{
			ID:        3,
			OwnerID:   1,
			FileName:  "file3",
			Metadata:  "md3",
			CreatedAt: time.Time{},
		},
	}

	for _, file := range files {
		_, err = s.AddFile(ctx, file)
		require.NoError(t, err)
	}

	res1, err := s.GetFiles(ctx, 1)
	assert.Len(t, res1, 0)
	assert.NoError(t, err)

	for _, file := range files {
		err = s.MarkFileAsUploaded(ctx, file.ID, file.OwnerID)
		require.NoError(t, err)
	}

	res2, err := s.GetFiles(ctx, 1)
	for _, r := range res2 {
		r.CreatedAt = time.Time{}
	}
	assert.NoError(t, err)
	assert.ElementsMatch(t, files, res2)
}

func TestStorage_CreateUser(t *testing.T) {

	ctx := context.Background()

	cfg, err := parseTestConfig()
	require.NoError(t, err)

	s, err := setupStorage(cfg, false)
	defer cleanup(cfg)
	require.NoError(t, err)

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

	s, err := setupStorage(cfg, false)
	defer cleanup(cfg)
	require.NoError(t, err)

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
