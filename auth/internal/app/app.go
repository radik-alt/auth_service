package app

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/radik-alt/auth_service/auth/internal/config"
	pb "github.com/radik-alt/auth_service/auth/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type App struct {
	cfg        *config.Config
	provider   *serviceProvider
	grpcServer *grpc.Server
}

func NewApp(ctx context.Context, envPath string) (*App, error) {
	cfg, err := config.Load(envPath)
	if err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}

	a := &App{cfg: cfg, provider: newServiceProvider(cfg)}

	authHandler, err := a.provider.AuthHandler(ctx)
	if err != nil {
		return nil, fmt.Errorf("init auth handler: %w", err)
	}

	a.grpcServer = grpc.NewServer()
	pb.RegisterAuthServiceServer(a.grpcServer, authHandler)
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
