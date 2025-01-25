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
	// create configuration manager
	cm, err := configs.ConfigManager()
	if err != nil {
		os.Exit(1)
	}

	// get current config
	config := cm.Get()

	// setup routes
	mux := http.NewServeMux()
	handlers.RegisterFactoryDefaultRoute(mux)
	handlers.RegisterRestartRoute(mux)

	// serve ui
	server := http.FileServer(http.FS(ui.DistDirFS))
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		server.ServeHTTP(w, r)

	})

	// listen
	log.Fatal(http.ListenAndServe(config.Server.HTTP.Port, mux))
}
