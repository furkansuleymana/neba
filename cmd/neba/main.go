package main

import (
	"io/fs"
	"log"
	"net/http"

	"github.com/furkansuleymana/neba/ui"
)

func main() {

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
