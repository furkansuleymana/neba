package handlers

import (
	"net/http"

	"github.com/furkansuleymana/neba/ui"
)

// RegisterStaticRoutes sets up handlers for serving static files (CSS, JS, images, etc.)
func RegisterStaticRoutes(mux *http.ServeMux) {
	// Create file server from embedded filesystem
	fs := http.FileServer(http.FS(ui.TemplatesDirFS))

	// Register handler for static files
	mux.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
		// Set appropriate cache control headers for static assets
		w.Header().Set("Cache-Control", "max-age=86400") // Cache for 1 day

		// Serve the requested file
		fs.ServeHTTP(w, r)
	})
}
