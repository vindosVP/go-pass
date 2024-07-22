package grpcmd

import (
	"context"
	"errors"
	"strconv"

	"google.golang.org/grpc/metadata"

	"github.com/vindosVP/go-pass/pkg/logger/sl"
)

func ExtractUID(ctx context.Context) (int, error) {

	errFailedToGetUID := errors.New("failed to extract uid from metadata")

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		sl.Log.Info("no metadata in context")
		return 0, errFailedToGetUID
	}
	uidMd := md.Get("uid")
	if len(uidMd) != 1 {
		sl.Log.Info("no uid in metadata")
		return 0, errFailedToGetUID
	}
	uid, err := strconv.Atoi(uidMd[0])
	if err != nil {
		sl.Log.Error("failed to convert uid to int", err)
		return 0, errFailedToGetUID
	}
	return uid, nil
}
