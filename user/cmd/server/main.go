package main

import (
	"github.com/radik-alt/auth_service/user/internal/config"
	userv1 "github.com/radik-alt/auth_service/user/pkg/user/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
)

func main() {
	_, err := config.Load(".env")
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	grpcServer := grpc.NewServer()
	userv1.RegisterUserServiceServer(grpcServer, nil)
	reflection.Register(grpcServer)
}
