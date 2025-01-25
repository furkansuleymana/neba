package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/furkansuleymana/neba/configs"
	"github.com/furkansuleymana/neba/discovery"
	"github.com/furkansuleymana/neba/handlers"
	"github.com/furkansuleymana/neba/ui"
)

func main() {
	// Create configuration manager
	cm, err := configs.ConfigManager()
	if err != nil {
		os.Exit(1)
	}

	// Get current config
	config := cm.Get()

	// Setup routes
	mux := http.NewServeMux()
	handlers.RegisterFactoryDefaultRoute(mux)
	handlers.RegisterRestartRoute(mux)

	// Serve UI
	server := http.FileServer(http.FS(ui.DistDirFS))
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		server.ServeHTTP(w, r)
	})

	// TESTING
	deviceList, _ := discovery.DiscoverWithSSDP()
	for _, device := range deviceList {
		slog.Info("Found device", slog.String("device", device))
	}
	// TESTING

	// Go!
	log.Fatal(http.ListenAndServe(config.Server.HTTP.Port, mux))
}
