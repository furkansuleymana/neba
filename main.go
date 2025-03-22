package main

import (
	"log"
	"net/http"
	"os"

	"github.com/furkansuleymana/neba/configs"
	"github.com/furkansuleymana/neba/handlers"
	"github.com/pkg/browser"
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
	handlers.RegisterStaticRoutes(mux)
	handlers.RegisterIndexRoute(mux)
	handlers.RegisterDiscoverDevicesRoute(mux)

	// Open browser
	err = browser.OpenURL("http://" + config.Server.HTTP.Address + config.Server.HTTP.Port)
	if err != nil {
		log.Println("Failed to open browser:", err)
	}

	// Go!
	log.Println("Server starting on", config.Server.HTTP.Port)
	log.Fatal(http.ListenAndServe(config.Server.HTTP.Port, mux))
}
