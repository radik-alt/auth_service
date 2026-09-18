package main

import (
	"context"
	"log"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	authpb "github.com/radik-alt/auth_service/auth/pkg/auth/v1"
	"github.com/radik-alt/auth_service/gateway/internal/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	cfg, err := config.Load(".env")
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	ctx := context.Background()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	if err := authpb.RegisterAuthServiceHandlerFromEndpoint(ctx, mux, cfg.AuthGRPCAddr, opts); err != nil {
		log.Fatalf("register auth gateway: %v", err)
	}

	// TODO: когда user-сервис будет готов (рабочий handler + user.pb.gw.go),
	// зарегистрировать его тем же способом на этом mux.

	log.Printf("HTTP gateway started on %s", cfg.HTTPAddr)
	if err := http.ListenAndServe(cfg.HTTPAddr, mux); err != nil {
		log.Fatalf("http serve: %v", err)
	}
}
