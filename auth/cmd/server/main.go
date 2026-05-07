package main

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/radik-alt/auth_service/auth/internal/config"
	"github.com/radik-alt/auth_service/auth/internal/handler"
	pgrepo "github.com/radik-alt/auth_service/auth/internal/repository/postgres"
	"github.com/radik-alt/auth_service/auth/internal/service/auth"
	pb "github.com/radik-alt/auth_service/auth/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net"
	"net/http"
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

	lis, err := net.Listen("tcp", cfg.GRPCAddr)
	if err != nil {
		log.Fatalf("listen grpc: %v", err)
	}

	go func() {
		log.Printf("gRPC server started on %s", cfg.GRPCAddr)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("grpc serve: %v", err)
		}
	}()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	if err := pb.RegisterAuthServiceHandlerFromEndpoint(ctx, mux, cfg.GRPCAddr, opts); err != nil {
		log.Fatalf("register gateway: %v", err)
	}

	log.Printf("HTTP gateway started on %s", cfg.HTTPAddr)
	if err := http.ListenAndServe(cfg.HTTPAddr, mux); err != nil {
		log.Fatalf("http serve: %v", err)
	}
}
