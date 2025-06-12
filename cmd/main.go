package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"log/slog"

	"context"

	"github.com/liriquew/test_task/internal/app"
	"github.com/liriquew/test_task/internal/lib/config"
	"github.com/liriquew/test_task/pkg/logger"
	"github.com/liriquew/test_task/pkg/logger/sl"
)

func main() {
	cfg := config.MustLoad()
	log := logger.SetupPrettySlog("USER SERVICE")

	log.Info("loaded config", slog.Any("config", cfg))

	app := app.New(log, cfg)

	app.Run()

	c := make(chan os.Signal, 1)

	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	<-c

	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()
	if err := app.Close(ctx); err != nil {
		log.Warn("error while shutdown server", sl.Err(err))
	}
}
