package api

import (
	"log/slog"
)

func Restart() error {
	slog.Info("restart successful")
	return nil
}
