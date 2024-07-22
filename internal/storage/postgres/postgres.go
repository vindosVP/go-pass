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

// AddPassword adds a new login-password pair
func (s *Storage) AddPassword(ctx context.Context, pwd *models.Password) (int, error) {
	return retry.DoWithData(func() (int, error) {
		query := `insert into passwords (owner_id, login, password, metadata, created_at) 
					values ($1, $2, $3, $4, $5) returning id`
		row := s.db.QueryRow(ctx, query, pwd.OwnerID, pwd.Login, pwd.Password, pwd.Metadata, time.Now())
		var id int
		err := row.Scan(&id)
		if err != nil {
			return 0, err
		}
		return id, nil
	}, retryOpts()...)
}

// AddCard adds a new bank card
func (s *Storage) AddCard(ctx context.Context, card *models.Card) (int, error) {
	return retry.DoWithData(func() (int, error) {
		query := `insert into cards (owner_id, number, cvc, owner, date, metadata, created_at) 
					values ($1, $2, $3, $4, $5, $6, $7) returning id`
		row := s.db.QueryRow(ctx, query, card.OwnerID, card.Number, card.CVC, card.Owner, card.Date, card.Metadata, time.Now())
		var id int
		err := row.Scan(&id)
		if err != nil {
			return 0, err
		}
		return id, nil
	}, retryOpts()...)
}

// AddText adds a new text
func (s *Storage) AddText(ctx context.Context, t *models.Text) (int, error) {
	return retry.DoWithData(func() (int, error) {
		query := `insert into texts (owner_id, text, metadata, created_at) 
					values ($1, $2, $3, $4) returning id`
		row := s.db.QueryRow(ctx, query, t.OwnerID, t.Text, t.Metadata, time.Now())
		var id int
		err := row.Scan(&id)
		if err != nil {
			return 0, err
		}
		return id, nil
	}, retryOpts()...)
}

// AddFile adds a new file
func (s *Storage) AddFile(ctx context.Context, f *models.File) (int, error) {
	return retry.DoWithData(func() (int, error) {
		query := `insert into files (owner_id, filename, metadata, created_at) 
					values ($1, $2, $3, $4) returning id`
		row := s.db.QueryRow(ctx, query, f.OwnerID, f.FileName, f.Metadata, time.Now())
		var id int
		err := row.Scan(&id)
		if err != nil {
			return 0, err
		}
		return id, nil
	}, retryOpts()...)
}

// UpdatePassword updates the password
func (s *Storage) UpdatePassword(ctx context.Context, pwd *models.Password) error {
	return retry.Do(func() error {
		query := `update 
    				passwords 
				  set 
					login=$1, 
					password=$2, 
					metadata=$3
				  where
				    id = $4 and owner_id = $5`
		_, err := s.db.Exec(ctx, query, pwd.Login, pwd.Password, pwd.Metadata, pwd.ID, pwd.OwnerID)
		if err != nil {
			return err
		}
		return nil
	}, retryOpts()...)
}

// UpdateCard updates the bank card
func (s *Storage) UpdateCard(ctx context.Context, card *models.Card) error {
	return retry.Do(func() error {
		query := `update 
    				cards 
				  set 
					number=$1, 
					cvc=$2, 
					owner=$3,
					date=$4,
					metadata=$5
				  where
				    id = $6 and owner_id = $7`
		_, err := s.db.Exec(ctx, query, card.Number, card.CVC, card.Owner, card.Date, card.Metadata, card.ID, card.OwnerID)
		if err != nil {
			return err
		}
		return nil
	}, retryOpts()...)
}

// UpdateText updates the text
func (s *Storage) UpdateText(ctx context.Context, t *models.Text) error {
	return retry.Do(func() error {
		query := `update 
    				texts 
				  set 
					text=$1, 
					metadata=$2
				  where
				    id = $3 and owner_id = $4`
		_, err := s.db.Exec(ctx, query, t.Text, t.Metadata, t.ID, t.OwnerID)
		if err != nil {
			return err
		}
		return nil
	}, retryOpts()...)
}

// DeletePassword deletes the password
func (s *Storage) DeletePassword(ctx context.Context, id int, ownerID int) error {
	return retry.Do(func() error {
		query := `delete from  
    				passwords 
				  where
				    id = $1 and owner_id = $2`
		_, err := s.db.Exec(ctx, query, id, ownerID)
		if err != nil {
			return err
		}
		return nil
	}, retryOpts()...)
}

// DeleteCard deletes the bank card
func (s *Storage) DeleteCard(ctx context.Context, id int, ownerID int) error {
	return retry.Do(func() error {
		query := `delete from  
    				cards 
				  where
				    id = $1 and owner_id = $2`
		_, err := s.db.Exec(ctx, query, id, ownerID)
		if err != nil {
			return err
		}
		return nil
	}, retryOpts()...)
}

// DeleteText deletes the text
func (s *Storage) DeleteText(ctx context.Context, id int, ownerID int) error {
	return retry.Do(func() error {
		query := `delete from  
    				texts 
				  where
				    id = $1 and owner_id = $2`
		_, err := s.db.Exec(ctx, query, id, ownerID)
		if err != nil {
			return err
		}
		return nil
	}, retryOpts()...)
}

// DeleteFile deletes the file
func (s *Storage) DeleteFile(ctx context.Context, id int, ownerID int) error {
	return retry.Do(func() error {
		query := `delete from  
    				files 
				  where
				    id = $1 and owner_id = $2`
		_, err := s.db.Exec(ctx, query, id, ownerID)
		if err != nil {
			return err
		}
		return nil
	}, retryOpts()...)
}

// GetPasswords returns all passwords
func (s *Storage) GetPasswords(ctx context.Context, ownerID int) ([]*models.Password, error) {
	return retry.DoWithData(func() ([]*models.Password, error) {
		query := `select
    				id, owner_id, login, password, metadata, created_at
    			  from  
    				passwords 
				  where
				    owner_id = $1`
		rows, err := s.db.Query(ctx, query, ownerID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		pwds := make([]*models.Password, 0)
		for rows.Next() {
			pwd := &models.Password{}
			err = rows.Scan(&pwd.ID, &pwd.OwnerID, &pwd.Login, &pwd.Password, &pwd.Metadata, &pwd.CreatedAt)
			if err != nil {
				return nil, err
			}
			pwds = append(pwds, pwd)
		}
		return pwds, nil
	}, retryOpts()...)
}

// GetCards returns all bank cards
func (s *Storage) GetCards(ctx context.Context, ownerID int) ([]*models.Card, error) {
	return retry.DoWithData(func() ([]*models.Card, error) {
		query := `select
    				id, owner_id, number, cvc, owner, date, metadata, created_at
    			  from  
    				cards 
				  where
				    owner_id = $1`
		rows, err := s.db.Query(ctx, query, ownerID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		cards := make([]*models.Card, 0)
		for rows.Next() {
			c := &models.Card{}
			err = rows.Scan(&c.ID, &c.OwnerID, &c.Number, &c.CVC, &c.Owner, &c.Date, &c.Metadata, &c.CreatedAt)
			if err != nil {
				return nil, err
			}
			cards = append(cards, c)
		}
		return cards, nil
	}, retryOpts()...)
}

// GetTexts returns all texts
func (s *Storage) GetTexts(ctx context.Context, ownerID int) ([]*models.Text, error) {
	return retry.DoWithData(func() ([]*models.Text, error) {
		query := `select
    				id, owner_id, text, metadata, created_at
    			  from  
    				texts 
				  where
				    owner_id = $1`
		rows, err := s.db.Query(ctx, query, ownerID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		texts := make([]*models.Text, 0)
		for rows.Next() {
			t := &models.Text{}
			err = rows.Scan(&t.ID, &t.OwnerID, &t.Text, &t.Metadata, &t.CreatedAt)
			if err != nil {
				return nil, err
			}
			texts = append(texts, t)
		}
		return texts, nil
	}, retryOpts()...)
}

// GetFiles returns all files
func (s *Storage) GetFiles(ctx context.Context, ownerID int) ([]*models.File, error) {
	return retry.DoWithData(func() ([]*models.File, error) {
		query := `select
    				id, owner_id, filename, metadata, created_at
    			  from  
    				files 
				  where
				    owner_id = $1 and uploaded`
		rows, err := s.db.Query(ctx, query, ownerID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		files := make([]*models.File, 0)
		for rows.Next() {
			f := &models.File{}
			err = rows.Scan(&f.ID, &f.OwnerID, &f.FileName, &f.Metadata, &f.CreatedAt)
			if err != nil {
				return nil, err
			}
			files = append(files, f)
		}
		return files, nil
	}, retryOpts()...)
}

// GetFile returns a file.
func (s *Storage) GetFile(ctx context.Context, id int, ownerID int) (*models.File, error) {
	return retry.DoWithData(func() (*models.File, error) {
		query := `select
    				id, owner_id, filename, metadata, created_at
    			  from  
    				files 
				  where
				    id = $1 and owner_id = $2`
		row := s.db.QueryRow(ctx, query, id, ownerID)
		file := &models.File{}
		err := row.Scan(&file.ID, &file.OwnerID, &file.FileName, &file.Metadata, &file.CreatedAt)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, storage.ErrFileNotExist
			}
			return nil, err
		}
		return file, nil
	}, retryOpts()...)
}

// MarkFileAsUploaded marks the file as uploaded.
func (s *Storage) MarkFileAsUploaded(ctx context.Context, id int, ownerID int) error {
	return retry.Do(func() error {
		query := `update files set 
    				uploaded=true 
				  where
				    id = $1 and owner_id = $2`
		_, err := s.db.Exec(ctx, query, id, ownerID)
		if err != nil {
			return err
		}
		return nil
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
