package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/radik-alt/auth_service/user/internal/app"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	app, err := app.NewApp(ctx, ".env")
	if err != nil {
		return fmt.Errorf("init app: %w", err)
	}
	defer app.Close()

	return app.Run(ctx)
}
