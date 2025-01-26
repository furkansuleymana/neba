package main

import (
	"log"
	"net/http"
	"os"

	"github.com/furkansuleymana/neba/api"
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
	handlers.RegisterRestartRoute(mux)

	// Serve UI
	server := http.FileServer(http.FS(ui.DistDirFS))
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		server.ServeHTTP(w, r)
	})

	// TESTING
	err = api.Restart("192.168.33.207", "root", "pass")
	if err != nil {
		log.Fatal(err)
	}
	os.Exit(0)
	// TESTING

	// Go!
	log.Fatal(http.ListenAndServe(config.Server.HTTP.Port, mux))
}
