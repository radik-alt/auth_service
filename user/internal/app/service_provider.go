package app

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/radik-alt/auth_service/user/internal/config"
	"github.com/radik-alt/auth_service/user/internal/handler"
	"github.com/radik-alt/auth_service/user/internal/repository"
	"github.com/radik-alt/auth_service/user/internal/repository/postgres"
	"github.com/radik-alt/auth_service/user/internal/service"
	"github.com/radik-alt/auth_service/user/internal/service/user"
	userv1 "github.com/radik-alt/auth_service/user/pkg/user/v1"
)

type serviceProvider struct {
	cfg *config.AppConfig

	pgPool      *pgxpool.Pool
	userRepo    repository.UserRepository
	userService service.UserService
	userHandler userv1.UserServiceServer
}

func newServiceProvider(cfg *config.AppConfig) *serviceProvider {
	return &serviceProvider{cfg: cfg}
}

func (p *serviceProvider) PGPool(ctx context.Context) (*pgxpool.Pool, error) {
	if p.pgPool == nil {
		pool, err := postgres.NewPool(ctx, p.cfg.DB)
		if err != nil {
			return nil, err
		}
		p.pgPool = pool
	}
	return p.pgPool, nil
}

func (p *serviceProvider) UserRepo(ctx context.Context) (repository.UserRepository, error) {
	if p.userRepo == nil {
		pool, err := p.PGPool(ctx)
		if err != nil {
			return nil, err
		}
		p.userRepo = postgres.NewUserRepo(pool)
	}
	return p.userRepo, nil
}

func (p *serviceProvider) UserService(ctx context.Context) (service.UserService, error) {
	if p.userService == nil {
		userRepo, err := p.UserRepo(ctx)
		if err != nil {
			return nil, err
		}
		p.userService = user.NewUserService(userRepo)
	}
	return p.userService, nil
}

func (p *serviceProvider) UserHandler(ctx context.Context) (userv1.UserServiceServer, error) {
	if p.userHandler == nil {
		svc, err := p.UserService(ctx)
		if err != nil {
			return nil, err
		}
		p.userHandler = handler.NewUserHandler(svc)
	}
	return p.userHandler, nil
}

func (p *serviceProvider) Close() error {
	if p.pgPool != nil {
		p.pgPool.Close()
	}
	return nil
}
