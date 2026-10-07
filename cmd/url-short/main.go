package main

import (
	"log/slog"
	"os"
)

func main() {
	// todo INIT CONFID cleanenv

	cfg := config.Mustload()

	// loger: log/slog

	logger := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})

	// init storage : postgresql

	pg := pgx.

	// router: gin

}
