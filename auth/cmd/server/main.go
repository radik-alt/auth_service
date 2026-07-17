package main

import (
	"context"
	"github.com/radik-alt/auth_service/auth/internal/config"
	"github.com/radik-alt/auth_service/auth/internal/handler"
	pgrepo "github.com/radik-alt/auth_service/auth/internal/repository/postgres"
	"github.com/radik-alt/auth_service/auth/internal/service/auth"
	pb "github.com/radik-alt/auth_service/auth/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

func main() {
	cfg, err := config.Load(".env")
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	ctx := context.Background()

	dbPool, err := pgrepo.NewPool(ctx, cfg.DB)
	if err != nil {
		log.Fatalf("postgres: %v", err)
	}
	defer dbPool.Close()

	if err := pgrepo.Migrate(ctx, dbPool); err != nil {
		log.Fatalf("migrate: %v", err)
	}

	//redisClient, err := redisrepo.NewClient(cfg.RedisAddr)
	//if err != nil {
	//	log.Fatalf("redis: %v", err)
	//}
	//defer redisClient.Close()

	authRepository := pgrepo.NewRefreshTokenRepo(dbPool)
	service := auth.NewAuthService(authRepository, cfg)
	authHandler := handler.NewAuthHandler(service)

	grpcServer := grpc.NewServer()
	pb.RegisterAuthServiceServer(grpcServer, authHandler)
	reflection.Register(grpcServer)

	lis, err := net.Listen("tcp", cfg.GRPCAddr)
	if err != nil {
		log.Fatalf("listen grpc: %v", err)
	}

	log.Printf("gRPC server started on %s", cfg.GRPCAddr)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("grpc serve: %v", err)
	}
}
