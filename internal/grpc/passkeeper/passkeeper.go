package passkeepergrpc

import (
	"context"
	"errors"
	"log/slog"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/vindosVP/go-pass/internal/models"
	passkeeperv1 "github.com/vindosVP/go-pass/internal/proto/passkeeper"
	"github.com/vindosVP/go-pass/internal/services/passkeeper"
	"github.com/vindosVP/go-pass/pkg/grpcmd"
	"github.com/vindosVP/go-pass/pkg/logger/sl"
)

// Keeper represents the keeeper API.
//
//go:generate go run github.com/vektra/mockery/v2@v2.43.2 --name=Keeper
type Keeper interface {
	Save(ctx context.Context, e *models.Entity) (int, error)
	Update(ctx context.Context, e *models.Entity) error
	Delete(ctx context.Context, id int, ownerID int, t models.EntityType) error
	List(ctx context.Context, ownerID int) ([]*models.Entity, error)
	SaveFile(str passkeeperv1.PassKeeper_UploadFileServer) error
	DownloadFile(id int, ownerID int, str passkeeperv1.PassKeeper_DownloadFileServer) error
}

type server struct {
	passkeeperv1.UnimplementedPassKeeperServer
	k Keeper
}

// Register registers the passkeeper service.
func Register(gRPCServer *grpc.Server, k Keeper) {
	passkeeperv1.RegisterPassKeeperServer(gRPCServer, &server{k: k})
}

// AddEntity adds the entity.
func (s server) AddEntity(ctx context.Context, in *passkeeperv1.AddEntityRequest) (*passkeeperv1.AddEntityResponse, error) {

	lg := sl.Log
	lg.Info("handling add entity request")

	uid, err := grpcmd.ExtractUID(ctx)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to extract uid: %v", err)
	}

	e := grpcToDTO(in.Entity, uid)
	id, err := s.k.Save(ctx, e)
	if err != nil {
		lg.Error("failed to save entity", err)
		return nil, status.Errorf(codes.Internal, "failed to save entity")
	}

	lg.Info("saved entity", slog.Int("id", id))
	return &passkeeperv1.AddEntityResponse{Id: int64(id)}, nil
}

// UpdateEntity updates the entity.
func (s server) UpdateEntity(ctx context.Context, in *passkeeperv1.UpdateEntityRequest) (*passkeeperv1.UpdateEntityResponse, error) {

	lg := sl.Log
	lg.Info("handling update entity request")

	uid, err := grpcmd.ExtractUID(ctx)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to extract uid: %v", err)
	}

	e := grpcToDTO(in.Entity, uid)
	e.ID = int(in.Id)
	err = s.k.Update(ctx, e)
	if err != nil {
		if errors.Is(err, passkeeper.ErrUnableToUpdateFile) {
			lg.Info("unable to update file")
			return nil, status.Errorf(codes.InvalidArgument, "unable to update file")
		}
		lg.Error("failed to update entity", err)
		return nil, status.Errorf(codes.Internal, "failed to update entity")
	}

	lg.Info("updated entity", slog.Int("id", e.ID))
	return &passkeeperv1.UpdateEntityResponse{}, nil
}

// DeleteEntity deletes the entity.
func (s server) DeleteEntity(ctx context.Context, in *passkeeperv1.DeleteEntityRequest) (*passkeeperv1.DeleteEntityResponse, error) {

	lg := sl.Log
	lg.Info("handling delete entity request")

	uid, err := grpcmd.ExtractUID(ctx)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to extract uid: %v", err)
	}

	err = s.k.Delete(ctx, int(in.Id), uid, totype(in.Type))
	if err != nil {
		lg.Error("failed to delete entity", err)
		return nil, status.Errorf(codes.Internal, "failed to delete entity")
	}
	lg.Info("deleted entity", slog.Int("id", int(in.Id)))
	return &passkeeperv1.DeleteEntityResponse{}, nil
}

// ListEntities lists all entities.
func (s server) ListEntities(ctx context.Context, _ *passkeeperv1.ListEntitiesRequest) (*passkeeperv1.ListEntitiesResponse, error) {

	lg := sl.Log
	lg.Info("handling list entities request")

	uid, err := grpcmd.ExtractUID(ctx)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to extract uid: %v", err)
	}

	res, err := s.k.List(ctx, uid)
	if err != nil {
		lg.Error("failed to list entities", err)
		return nil, status.Errorf(codes.Internal, "failed to list entities")
	}
	resp := &passkeeperv1.ListEntitiesResponse{Entity: make([]*passkeeperv1.Entity, 0, len(res))}
	for _, e := range res {
		resp.Entity = append(resp.Entity, dtoToGRPC(e))
	}
	return resp, nil
}

// UploadFile uploads files to the server.
func (s server) UploadFile(str passkeeperv1.PassKeeper_UploadFileServer) error {
	return s.k.SaveFile(str)
}

// DownloadFile downloads file from the server.
func (s server) DownloadFile(in *passkeeperv1.DownloadFileRequest, str passkeeperv1.PassKeeper_DownloadFileServer) error {
	uid, err := grpcmd.ExtractUID(str.Context())
	if err != nil {
		return status.Errorf(codes.InvalidArgument, "failed to extract uid: %v", err)
	}
	return s.k.DownloadFile(int(in.Id), uid, str)
}

func dtoToGRPC(e *models.Entity) *passkeeperv1.Entity {

	var t passkeeperv1.Type
	switch e.Type {
	case models.TypePassword:
		t = passkeeperv1.Type_PASSWORD
	case models.TypeFile:
		t = passkeeperv1.Type_FILE
	case models.TypeCard:
		t = passkeeperv1.Type_CARD
	case models.TypeText:
		t = passkeeperv1.Type_TEXT
	}

	return &passkeeperv1.Entity{
		Id:         int64(e.ID),
		Type:       t,
		Login:      e.Login,
		Password:   e.Password,
		CardNumber: e.CardNumber,
		CardOwner:  e.CardOwner,
		CardCVC:    e.CardCVC,
		CardExp:    e.CardExp,
		Text:       e.Text,
		Filename:   e.Filename,
		Metadata:   e.Metadata,
	}
}

func grpcToDTO(e *passkeeperv1.Entity, ownerID int) *models.Entity {

	var t models.EntityType
	switch e.Type {
	case passkeeperv1.Type_PASSWORD:
		t = models.TypePassword
	case passkeeperv1.Type_CARD:
		t = models.TypeCard
	case passkeeperv1.Type_TEXT:
		t = models.TypeText
	case passkeeperv1.Type_FILE:
		t = models.TypeFile
	}

	return &models.Entity{
		ID:         int(e.Id),
		OwnerID:    ownerID,
		Type:       t,
		Login:      e.Login,
		Password:   e.Password,
		CardNumber: e.CardNumber,
		CardOwner:  e.CardOwner,
		CardCVC:    e.CardCVC,
		CardExp:    e.CardExp,
		Text:       e.Text,
		Filename:   e.Filename,
		Metadata:   e.Metadata,
	}
}

func totype(grpcType passkeeperv1.Type) models.EntityType {
	switch grpcType {
	case passkeeperv1.Type_PASSWORD:
		return models.TypePassword
	case passkeeperv1.Type_CARD:
		return models.TypeCard
	case passkeeperv1.Type_TEXT:
		return models.TypeText
	default:
		return models.TypeFile
	}
}
