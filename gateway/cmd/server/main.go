package main

import (
	"context"
	"log"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	authpb "github.com/radik-alt/auth_service/auth/pkg/auth/v1"
	"github.com/radik-alt/auth_service/gateway/internal/client"
	"github.com/radik-alt/auth_service/gateway/internal/config"
	"github.com/radik-alt/auth_service/gateway/internal/middleware"
	userv1 "github.com/radik-alt/auth_service/user/pkg/user/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	cfg, err := config.Load(".env")
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	ctx := context.Background()

	mux := runtime.NewServeMux(
		runtime.WithIncomingHeaderMatcher(middleware.BlockIncomingHeaders),
		runtime.WithMetadata(middleware.UserIDMetadata),
	)
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	if err := authpb.RegisterAuthServiceHandlerFromEndpoint(ctx, mux, cfg.AuthGRPCAddr, opts); err != nil {
		log.Fatalf("register auth gateway: %v", err)
	}
	if err := userv1.RegisterUserServiceHandlerFromEndpoint(ctx, mux, cfg.UserGRPCAddr, opts); err != nil {
		log.Fatalf("register user gateway: %v", err)
	}

	// отдельное соединение для проверки токенов в middleware; подключается при первом вызове
	authConn, err := grpc.NewClient(cfg.AuthGRPCAddr, opts...)
	if err != nil {
		log.Fatalf("auth grpc client: %v", err)
	}
	defer authConn.Close()

	handler := middleware.Auth(client.NewAuthClient(authConn))(mux)

	log.Printf("HTTP gateway started on %s", cfg.HTTPAddr)
	if err := http.ListenAndServe(cfg.HTTPAddr, handler); err != nil {
		log.Fatalf("http serve: %v", err)
	}
}
