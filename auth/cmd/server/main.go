package main

import (
	"context"
	"log"
	"net"

	"github.com/radik-alt/auth_service/auth/internal/config"
	"github.com/radik-alt/auth_service/auth/internal/handler"
	redisrepo "github.com/radik-alt/auth_service/auth/internal/repository/redis"
	"github.com/radik-alt/auth_service/auth/internal/service/auth"
	pb "github.com/radik-alt/auth_service/auth/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	ctx := context.Background()

	cfg, err := config.Load(".env")
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	redisClient, err := redisrepo.NewClient(ctx, cfg.RedisAddr)
	if err != nil {
		log.Fatalf("redis: %v", err)
	}
	defer redisClient.Close()

	refreshRepo := redisrepo.NewRefreshTokenRepo(redisClient)
	blacklist := redisrepo.NewBlacklistRepo(redisClient)

	authService := auth.NewAuthService(refreshRepo, blacklist, cfg)
	authHandler := handler.NewAuthHandler(authService)

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
