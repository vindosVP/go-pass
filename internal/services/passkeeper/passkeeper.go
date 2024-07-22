package passkeeper

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/vindosVP/go-pass/internal/filemanager"
	"github.com/vindosVP/go-pass/internal/models"
	passkeeperv1 "github.com/vindosVP/go-pass/internal/proto/passkeeper"
	"github.com/vindosVP/go-pass/internal/storage"
	"github.com/vindosVP/go-pass/pkg/grpcmd"
	"github.com/vindosVP/go-pass/pkg/logger/sl"
)

var (

	// ErrUnknownEntity - error if entity`s type is unknown
	ErrUnknownEntity = errors.New("unknown entity")

	// ErrUnableToUpdateFile - error if tried to update file, to update delete file and add another one
	ErrUnableToUpdateFile = errors.New("unable to update file")

	// ErrUnableToSaveFile - error if tried to save file, to save file, use SaveFile method
	ErrUnableToSaveFile = errors.New("unable to save file")
)

// PasswordStorage is a password storage API
//
//go:generate go run github.com/vektra/mockery/v2@v2.43.2 --name=PasswordStorage
type PasswordStorage interface {
	AddPassword(ctx context.Context, pwd *models.Password) (int, error)
	UpdatePassword(ctx context.Context, pwd *models.Password) error
	DeletePassword(ctx context.Context, id int, ownerID int) error
	GetPasswords(ctx context.Context, ownerID int) ([]*models.Password, error)
}

// CardStorage is a card storage API
//
//go:generate go run github.com/vektra/mockery/v2@v2.43.2 --name=CardStorage
type CardStorage interface {
	AddCard(ctx context.Context, card *models.Card) (int, error)
	UpdateCard(ctx context.Context, card *models.Card) error
	DeleteCard(ctx context.Context, id int, ownerID int) error
	GetCards(ctx context.Context, ownerID int) ([]*models.Card, error)
}

// TextStorage is a text storage API
//
//go:generate go run github.com/vektra/mockery/v2@v2.43.2 --name=TextStorage
type TextStorage interface {
	AddText(ctx context.Context, t *models.Text) (int, error)
	UpdateText(ctx context.Context, t *models.Text) error
	DeleteText(ctx context.Context, id int, ownerID int) error
	GetTexts(ctx context.Context, ownerID int) ([]*models.Text, error)
}

// FileStorage is a file storage API
//
//go:generate go run github.com/vektra/mockery/v2@v2.43.2 --name=FileStorage
type FileStorage interface {
	GetFiles(ctx context.Context, ownerID int) ([]*models.File, error)
	GetFile(ctx context.Context, id int, ownerID int) (*models.File, error)
	DeleteFile(ctx context.Context, id int, ownerID int) error
	AddFile(ctx context.Context, f *models.File) (int, error)
	MarkFileAsUploaded(ctx context.Context, id int, ownerID int) error
}

type Keeper struct {
	ps    PasswordStorage
	cs    CardStorage
	ts    TextStorage
	fs    FileStorage
	fPath string
}

const chunkSize = 4 * 1024

// List returns all user`s entities.
func (k *Keeper) List(ctx context.Context, ownerID int) ([]*models.Entity, error) {

	sl.Log.Info("getting list of entities")

	res := make([]*models.Entity, 0)

	pwds, err := k.ps.GetPasswords(ctx, ownerID)
	if err != nil {
		return nil, err
	}
	for _, pwd := range pwds {
		res = append(res, pwd.ToEntity())
	}

	cards, err := k.cs.GetCards(ctx, ownerID)
	if err != nil {
		return nil, err
	}
	for _, card := range cards {
		res = append(res, card.ToEntity())
	}

	texts, err := k.ts.GetTexts(ctx, ownerID)
	if err != nil {
		return nil, err
	}
	for _, text := range texts {
		res = append(res, text.ToEntity())
	}

	files, err := k.fs.GetFiles(ctx, ownerID)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		res = append(res, file.ToEntity())
	}

	return res, nil
}

// Save saves the entity (password, card or text).
func (k *Keeper) Save(ctx context.Context, e *models.Entity) (int, error) {
	switch e.Type {
	case models.TypePassword:
		sl.Log.Info("adding new password")
		return k.ps.AddPassword(ctx, e.ToPassword())
	case models.TypeCard:
		sl.Log.Info("adding new card")
		return k.cs.AddCard(ctx, e.ToCard())
	case models.TypeText:
		sl.Log.Info("adding new text")
		return k.ts.AddText(ctx, e.ToText())
	case models.TypeFile:
		sl.Log.Error("unable to save file")
		return 0, ErrUnableToSaveFile
	}
	sl.Log.Error("unknown entity type", slog.String("type", string(e.Type)))
	return 0, ErrUnknownEntity
}

// Delete deletes the entity.
func (k *Keeper) Delete(ctx context.Context, id int, ownerID int, t models.EntityType) error {
	switch t {
	case models.TypePassword:
		sl.Log.Info("deleting password", slog.Int("id", id))
		return k.ps.DeletePassword(ctx, id, ownerID)
	case models.TypeCard:
		sl.Log.Info("deleting card", slog.Int("id", id))
		return k.cs.DeleteCard(ctx, id, ownerID)
	case models.TypeText:
		sl.Log.Info("deleting text", slog.Int("id", id))
		return k.ts.DeleteText(ctx, id, ownerID)
	case models.TypeFile:
		sl.Log.Info("deleting text", slog.Int("id", id))
		file, err := k.fs.GetFile(ctx, id, ownerID)
		if err != nil {
			return err
		}
		deleter := filemanager.NewFileDeleter()
		filename := fmt.Sprintf("%d_%s", file.ID, file.FileName)
		deleter.SetFile(filename, k.fPath)
		err = deleter.Delete()
		if err != nil {
			return err
		}
		return k.fs.DeleteFile(ctx, id, ownerID)
	}
	sl.Log.Error("unknown entity type", slog.String("type", string(t)))
	return ErrUnknownEntity
}

// Update updates the entity (password, card or text).
func (k *Keeper) Update(ctx context.Context, e *models.Entity) error {
	switch e.Type {
	case models.TypePassword:
		sl.Log.Info("updating password", slog.Int("id", e.ID))
		return k.ps.UpdatePassword(ctx, e.ToPassword())
	case models.TypeCard:
		sl.Log.Info("updating card", slog.Int("id", e.ID))
		return k.cs.UpdateCard(ctx, e.ToCard())
	case models.TypeText:
		sl.Log.Info("updating text", slog.Int("id", e.ID))
		return k.ts.UpdateText(ctx, e.ToText())
	case models.TypeFile:
		sl.Log.Error("unable to update file")
		return ErrUnableToUpdateFile
	}
	sl.Log.Error("unknown entity type", slog.String("type", string(e.Type)))
	return ErrUnknownEntity
}

// SaveFile saves the file.
func (k *Keeper) SaveFile(str passkeeperv1.PassKeeper_UploadFileServer) error {

	lg := sl.Log
	lg.Info("handling upload file request")

	savedToDB := false
	var fileId int
	var ownerID int

	f := filemanager.NewFileSaver()
	defer f.Close()

	for {
		req, err := str.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			lg.Error("failed to save file", sl.Err(err))
			return status.Errorf(codes.Internal, "failed to upload file")
		}
		if !savedToDB {
			uid, err := grpcmd.ExtractUID(str.Context())
			if err != nil {
				return status.Errorf(codes.InvalidArgument, "failed to extract uid")
			}
			ownerID = uid
			file := &models.File{
				OwnerID:  uid,
				FileName: req.Filename,
				Metadata: req.Metadata,
			}
			id, err := k.fs.AddFile(str.Context(), file)
			if err != nil {
				lg.Error("failed to save file", sl.Err(err))
				return status.Errorf(codes.Internal, "failed to save file")
			}
			fileId = id
			savedToDB = true
		}
		if !f.IsFileSet() {
			filename := fmt.Sprintf("%d_%s", fileId, req.GetFilename())
			err := f.SetFile(filename, k.fPath)
			if err != nil {
				lg.Error("failed to save file", sl.Err(err))
				return status.Errorf(codes.Internal, "failed to save file")
			}
		}

		ch := req.GetChunk()
		if err = f.Write(ch); err != nil {
			lg.Error("failed to save file", sl.Err(err))
			return status.Errorf(codes.Internal, "failed to upload file")
		}
	}

	err := k.fs.MarkFileAsUploaded(str.Context(), fileId, ownerID)
	if err != nil {
		lg.Error("failed to save file", sl.Err(err))
		return status.Errorf(codes.Internal, "failed to upload file")
	}

	err = str.SendAndClose(&passkeeperv1.UploadFileResponse{Id: int64(fileId)})
	if err != nil {
		lg.Error("failed to save file", sl.Err(err))
		return status.Errorf(codes.Internal, "failed to save file")
	}

	lg.Info("saved file")
	return nil

}

// DownloadFile returns the file to download.
func (k *Keeper) DownloadFile(id int, ownerID int, str passkeeperv1.PassKeeper_DownloadFileServer) error {

	lg := sl.Log
	lg.Info("handling download file request")

	file, err := k.fs.GetFile(str.Context(), id, ownerID)
	if err != nil {
		if errors.Is(err, storage.ErrFileNotExist) {
			lg.Info("file not found")
			return status.Errorf(codes.NotFound, "file not found")
		}
		lg.Error("failed to download file", sl.Err(err))
		return status.Errorf(codes.Internal, "failed to download file")
	}

	filename := fmt.Sprintf("%d_%s", file.ID, file.FileName)
	fileReader := filemanager.NewFileReader(chunkSize)
	err = fileReader.SetFile(filename, k.fPath)
	if err != nil {
		lg.Error("failed to download file", sl.Err(err))
		return status.Errorf(codes.Internal, "failed to download file")
	}
	defer fileReader.Close()
	for fileReader.Next() {
		resp := &passkeeperv1.DownloadFileResponse{
			Filename: file.FileName,
			Chunk:    fileReader.Data(),
		}
		err = str.Send(resp)
		if err != nil {
			lg.Error("failed to download file", sl.Err(err))
			return status.Errorf(codes.Internal, "failed to download file")
		}
	}

	return nil
}

// New creates a new Keeper instance
func New(ps PasswordStorage, cs CardStorage, ts TextStorage, fs FileStorage, fPath string) *Keeper {
	return &Keeper{ps: ps, cs: cs, ts: ts, fs: fs, fPath: fPath}
}
