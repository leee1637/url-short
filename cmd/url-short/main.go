package main

import (
	"log/slog"
	"os"
)

func main() {
	// todo INIT CONFID cleanenv

	// loger: log/slog

	logger := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})

	// init storage : postgresql

	// router: gin

	// loger: log/slog
}
