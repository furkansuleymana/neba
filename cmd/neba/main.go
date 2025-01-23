package main

import (
	"io/fs"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/furkansuleymana/neba/settings"
	"github.com/furkansuleymana/neba/ui"
)

func main() {
	// create settings manager
	cm, err := settings.SettingsManager()
	if err != nil {
		os.Exit(1)
	}

	// get current settings
	setting := cm.Get()
	slog.Info("retrieved settings", slog.Any("settings", setting))

	// update settings
	err = cm.Update(func(set *settings.ApplicationSettings) {
		set.Server.HTTP.Address = "127.0.0.1"
		set.Server.HTTP.Port = "8080"
		set.Database.Path = "/new/path/to/database"
	})
	if err != nil {
		slog.Error("failed to update settings", "error", err)
	}

	// get updated settings
	setting = cm.Get()
	slog.Info("retrieved settings", slog.Any("settings", setting))

	//

	fileSystem, err := fs.Sub(ui.Dist, "dist")
	if err != nil {
		log.Fatalf("failed to create filesystem: %v", err)
	}

	http.Handle("/", http.FileServer(http.FS(fileSystem)))

	addr := ":8080"
	log.Printf("starting server on %s\n", addr)

	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}

}
