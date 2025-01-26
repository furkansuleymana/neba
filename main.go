package main

import (
	"log"
	"net/http"
	"os"

	"github.com/furkansuleymana/neba/configs"
	"github.com/furkansuleymana/neba/database"
	"github.com/furkansuleymana/neba/database/models"
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
	device := models.AxisDevice{
		SerialNumber: "A12B34C56D78",
		IPAddress:    "192.168.1.1",
		Model:        "M1075-L",
		OSVersion:    "11.2.68",
	}
	bucketName := "devices"

	db, err := database.Open(config.Database.Path, bucketName)
	if err != nil {
		log.Fatalf("Could not initialize database: %v", err)
	}
	defer database.CloseDB(db)

	if err := database.Update(db, bucketName, device); err != nil {
		log.Fatalf("Could not save device: %v", err)
	}

	retrievedDevice, err := database.View(db, bucketName, "A12B34C56D78")
	if err != nil {
		log.Fatalf("Could not get device: %v", err)
	}

	log.Println(retrievedDevice)
	os.Exit(0)
	// TESTING

	// Go!
	log.Fatal(http.ListenAndServe(config.Server.HTTP.Port, mux))
}
