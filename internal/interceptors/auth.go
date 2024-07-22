package interceptors

import (
	"context"
	"strconv"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/vindosVP/go-pass/internal/jwt"
	"github.com/vindosVP/go-pass/pkg/logger/sl"
)

type AuthInterceptor struct {
	secret string
}

func NewAuthInterceptor(secret string) *AuthInterceptor {
	return &AuthInterceptor{secret: secret}
}

var openedMethods = []string{"/auth.Auth/Login", "/auth.Auth/Register"}

type wrappedStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (w *wrappedStream) Context() context.Context {
	return w.ctx
}

func (w *wrappedStream) SetContext(ctx context.Context) {
	w.ctx = ctx
}

func (w *wrappedStream) RecvMsg(m interface{}) error {
	return w.ServerStream.RecvMsg(m)
}

func (w *wrappedStream) SendMsg(m interface{}) error {
	return w.ServerStream.SendMsg(m)
}

type StreamContextWrapper interface {
	grpc.ServerStream
	SetContext(context.Context)
}

func (a *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {

		lg := sl.Log

		for _, v := range openedMethods {
			if v == info.FullMethod {
				return handler(ctx, req)
			}
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			lg.Info("missing token")
			return nil, status.Errorf(codes.Unauthenticated, "missing token")
		}
		t := md.Get("token")
		if len(t) != 1 {
			lg.Info("wrong token format")
			return nil, status.Errorf(codes.Unauthenticated, "wrong token format")
		}
		token := t[0]

		email, uid, err := jwt.VerifyToken(token, a.secret)
		if err != nil {
			lg.Info("invalid token")
			return nil, status.Errorf(codes.Unauthenticated, "invalid token")
		}

		newMD := map[string]string{
			"email": email,
			"uid":   strconv.Itoa(uid),
		}

		return handler(metadata.NewIncomingContext(ctx, metadata.New(newMD)), req)
	}
}

func (a *AuthInterceptor) Stream() grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		stream grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {

		lg := sl.Log

		for _, v := range openedMethods {
			if v == info.FullMethod {
				return handler(srv, stream)
			}
		}

		md, ok := metadata.FromIncomingContext(stream.Context())
		if !ok {
			lg.Info("missing token")
			return status.Errorf(codes.Unauthenticated, "missing token")
		}
		t := md.Get("token")
		if len(t) != 1 {
			lg.Info("wrong token format")
			return status.Errorf(codes.Unauthenticated, "wrong token format")
		}
		token := t[0]

		email, uid, err := jwt.VerifyToken(token, a.secret)
		if err != nil {
			lg.Info("invalid token")
			return status.Errorf(codes.Unauthenticated, "invalid token")
		}

		newMD := map[string]string{
			"email": email,
			"uid":   strconv.Itoa(uid),
		}
		newCtx := metadata.NewIncomingContext(stream.Context(), metadata.New(newMD))

		return handler(srv, &wrappedStream{stream, newCtx})
	}
}
