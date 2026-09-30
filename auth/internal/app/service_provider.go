package app

import (
	"context"
	"errors"
	"fmt"

	"github.com/radik-alt/auth_service/auth/internal/client"
	userclient "github.com/radik-alt/auth_service/auth/internal/client/user"
	"github.com/radik-alt/auth_service/auth/internal/config"
	"github.com/radik-alt/auth_service/auth/internal/handler"
	"github.com/radik-alt/auth_service/auth/internal/repository"
	redisrepo "github.com/radik-alt/auth_service/auth/internal/repository/redis"
	"github.com/radik-alt/auth_service/auth/internal/service"
	"github.com/radik-alt/auth_service/auth/internal/service/auth"
	pb "github.com/radik-alt/auth_service/auth/pkg/auth/v1"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type serviceProvider struct {
	cfg *config.Config

	redisClient *redis.Client
	refreshRepo repository.TokenRepository
	blacklist   repository.TokenBlacklist
	userConn    *grpc.ClientConn
	userClient  client.UserClient
	authService service.AuthService
	authHandler pb.AuthServiceServer
}

func newServiceProvider(cfg *config.Config) *serviceProvider {
	return &serviceProvider{cfg: cfg}
}

func (p *serviceProvider) RedisClient(ctx context.Context) (*redis.Client, error) {
	if p.redisClient == nil {
		client, err := redisrepo.NewClient(ctx, p.cfg.RedisAddr)
		if err != nil {
			return nil, err
		}
		p.redisClient = client
	}
	return p.redisClient, nil
}

func (p *serviceProvider) RefreshRepo(ctx context.Context) (repository.TokenRepository, error) {
	if p.refreshRepo == nil {
		client, err := p.RedisClient(ctx)
		if err != nil {
			return nil, err
		}
		p.refreshRepo = redisrepo.NewRefreshTokenRepo(client)
	}
	return p.refreshRepo, nil
}

func (p *serviceProvider) Blacklist(ctx context.Context) (repository.TokenBlacklist, error) {
	if p.blacklist == nil {
		client, err := p.RedisClient(ctx)
		if err != nil {
			return nil, err
		}

		p.blacklist = redisrepo.NewBlacklistRepo(client)
	}

	return p.blacklist, nil
}

// UserConn не подключается сразу: grpc.NewClient соединяется при первом вызове,
// поэтому auth стартует, даже если user-сервис ещё не поднялся.
func (p *serviceProvider) UserConn() (*grpc.ClientConn, error) {
	if p.userConn == nil {
		conn, err := grpc.NewClient(p.cfg.UserGRPCAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return nil, fmt.Errorf("create user grpc client: %w", err)
		}
		p.userConn = conn
	}
	return p.userConn, nil
}

func (p *serviceProvider) UserClient() (client.UserClient, error) {
	if p.userClient == nil {
		conn, err := p.UserConn()
		if err != nil {
			return nil, err
		}
		p.userClient = userclient.NewUserClient(conn)
	}
	return p.userClient, nil
}

func (p *serviceProvider) AuthService(ctx context.Context) (service.AuthService, error) {
	if p.authService == nil {
		refreshRepo, err := p.RefreshRepo(ctx)
		if err != nil {
			return nil, err
		}
		blacklist, err := p.Blacklist(ctx)
		if err != nil {
			return nil, err
		}
		userClient, err := p.UserClient()
		if err != nil {
			return nil, err
		}
		p.authService = auth.NewAuthService(refreshRepo, blacklist, userClient, &config.Config{
			JWTSecret:  p.cfg.JWTSecret,
			AccessTTL:  p.cfg.AccessTTL,
			RefreshTTL: p.cfg.RefreshTTL,
		})
	}
	return p.authService, nil
}

func (p *serviceProvider) AuthHandler(ctx context.Context) (pb.AuthServiceServer, error) {
	if p.authHandler == nil {
		svc, err := p.AuthService(ctx)
		if err != nil {
			return nil, err
		}
		p.authHandler = handler.NewAuthHandler(svc)
	}
	return p.authHandler, nil
}

func (p *serviceProvider) Close() error {
	var errs []error
	if p.userConn != nil {
		errs = append(errs, p.userConn.Close())
	}
	if p.redisClient != nil {
		errs = append(errs, p.redisClient.Close())
	}
	return errors.Join(errs...)
}
