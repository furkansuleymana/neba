package main

import (
	"log"
	"net/http"
	"os"

	"github.com/furkansuleymana/neba/configs"
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
	handlers.RegisterFindDevicesRoute(mux)
	handlers.RegisterRestartRoute(mux)
	handlers.RegisterStreamRoute(mux)
	handlers.RegisterDevicesListRoute(mux)

	// Serve UI
	server := http.FileServer(http.FS(ui.TemplatesDirFS))
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		server.ServeHTTP(w, r)
	})

	// Go!
	log.Fatal(http.ListenAndServe(config.Server.HTTP.Port, mux))
}
