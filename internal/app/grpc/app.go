package grpcapp

import (
	"context"
	"fmt"
	"log/slog"
	"net"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	authgrpc "github.com/vindosVP/go-pass/internal/grpc/auth"
	"github.com/vindosVP/go-pass/pkg/logger/sl"
)

// App is a grpc app representation.
type App struct {
	grpcServer *grpc.Server
	port       int
}

// MustRun runs gRPC server and panics if any error occurs.
func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

// Run runs gRPC server.
func (a *App) Run() error {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}
	sl.Log.Info("grpc server started", slog.String("addr", l.Addr().String()))
	if err := a.grpcServer.Serve(l); err != nil {
		return fmt.Errorf("failed to serve: %w", err)
	}
	return nil
}

// Stop stops gRPC server.
func (a *App) Stop() {
	sl.Log.Info("stopping gRPC server", slog.Int("port", a.port))
	a.grpcServer.GracefulStop()
}

// New creates a grpc app instance.
func New(port int, auth authgrpc.Auth) *App {
	loggingOpts := []logging.Option{
		logging.WithLogOnEvents(
			logging.PayloadReceived, logging.PayloadSent,
		),
	}
	recoveryOpts := []recovery.Option{
		recovery.WithRecoveryHandler(func(p interface{}) (err error) {
			sl.Log.Error("recovered from panic", slog.Any("panic", p))
			return status.Errorf(codes.Internal, "internal error")
		}),
	}
	grpcServer := grpc.NewServer(grpc.ChainUnaryInterceptor(
		recovery.UnaryServerInterceptor(recoveryOpts...),
		logging.UnaryServerInterceptor(InterceptorLogger(sl.Log), loggingOpts...),
	))

	authgrpc.Register(grpcServer, auth)

	return &App{
		grpcServer: grpcServer,
		port:       port,
	}
}

// InterceptorLogger adapts slog logger to interceptor logger.
func InterceptorLogger(l *slog.Logger) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		l.Log(ctx, slog.Level(lvl), msg, fields...)
	})
}
