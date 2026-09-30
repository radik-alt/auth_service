package app

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/radik-alt/auth_service/user/internal/config"
	userv1 "github.com/radik-alt/auth_service/user/pkg/user/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type App struct {
	cfg        *config.AppConfig
	provider   *serviceProvider
	grpcServer *grpc.Server
}

func NewApp(ctx context.Context, envPath string) (*App, error) {
	cfg, err := config.Load(envPath)
	if err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}

	a := &App{cfg: cfg, provider: newServiceProvider(cfg)}

	userHandler, err := a.provider.UserHandler(ctx)
	if err != nil {
		a.provider.Close()
		return nil, fmt.Errorf("init user handler: %w", err)
	}

	a.grpcServer = grpc.NewServer()
	userv1.RegisterUserServiceServer(a.grpcServer, userHandler)
	reflection.Register(a.grpcServer)

	return a, nil
}

func (a *App) Run(ctx context.Context) error {
	lis, err := net.Listen("tcp", a.cfg.GRPCAddr)
	if err != nil {
		return fmt.Errorf("listen: %w", err)
	}

	errCh := make(chan error, 1)
	go func() { errCh <- a.grpcServer.Serve(lis) }()
	log.Printf("gRPC server started on %s", a.cfg.GRPCAddr)

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		a.grpcServer.GracefulStop()
		return nil
	}
}

func (a *App) Close() error { return a.provider.Close() }
