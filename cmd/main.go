package main

import (
	"log/slog"

	_ "github.com/ucpr/mongo-streamer/internal/log"
)

func main() {
	slog.Info("Hello World")
}
