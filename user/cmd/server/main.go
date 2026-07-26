package main

import (
	"github.com/radik-alt/auth_service/user/internal/config"
	proto "github.com/radik-alt/auth_service/user/proto"
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
	proto.RegisterUserServiceServer(grpcServer, nil)
	reflection.Register(grpcServer)
}
