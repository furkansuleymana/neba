package main

import (
	"log"
	"net/http"
	"os"

	"github.com/furkansuleymana/neba/configs"
	"github.com/furkansuleymana/neba/handlers"
	"github.com/furkansuleymana/neba/ui"
	"github.com/pkg/browser"
)

func main() {
	// Create configuration manager
	cm, err := configs.ConfigManager()
	if err != nil {
		log.Fatal("Failed to create config manager:", err)
		os.Exit(1)
	}

	// Get current config
	config := cm.Get()

	// Setup server
	fs := http.FileServer(http.FS(ui.FS))
	mux := http.NewServeMux()

	mux.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	})

	handlers.RegisterRootRoute(fs, mux)
	handlers.RegisterHomeRoute(fs, mux)
	handlers.RegisterFindDevicesRoute(fs, mux)
	handlers.RegisterManageDevicesRoute(fs, mux)

	// Open browser
	err = browser.OpenURL("http://" + config.Server.HTTP.Address + config.Server.HTTP.Port)
	if err != nil {
		log.Println("Failed to open browser:", err)
	}

	// Go!
	log.Fatal(http.ListenAndServe(config.Server.HTTP.Port, mux))
	log.Println("Hello from Neba!")
}
