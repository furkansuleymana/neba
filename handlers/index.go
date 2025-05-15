package handlers

import (
	"html/template"
	"net/http"
	"strings"

	"github.com/furkansuleymana/neba/ui"
)

func RegisterIndexRoute(mux *http.ServeMux) {
	// Create file server for static files
	fs := http.FileServer(http.FS(ui.TemplatesDirFS))
	// Handle root path
	mux.HandleFunc("/", handleIndex)
	// Handle static files
	mux.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "max-age=86400") // Cache for 1 day
		fs.ServeHTTP(w, r)
	})
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	// If path is not root, and not a static file, let the file server handle it
	if r.URL.Path != "/" && !strings.HasPrefix(r.URL.Path, "/static/") {
		http.FileServer(http.FS(ui.TemplatesDirFS)).ServeHTTP(w, r)
		return
	}

	// Parse templates from embedded filesystem
	tmpl, err := template.ParseFS(ui.TemplatesDirFS,
		"index.html",
		"home.html")
	if err != nil {
		http.Error(w, "Failed to parse templates: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Execute template with page data
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Failed to render template: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
