package handlers

import (
	"html/template"
	"log"
	"net/http"

	"github.com/furkansuleymana/neba/ui"
)

var (
	rootTmpl *template.Template
)

func RegisterRootRoute(fs http.Handler, mux *http.ServeMux) {
	var err error
	rootTmpl, err = template.ParseFS(ui.FS, "index.html")
	if err != nil {
		log.Fatalf("Failed to parse templates: %v", err)
	}

	mux.HandleFunc("/", handleRoot)
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	if err := rootTmpl.ExecuteTemplate(w, "index.html", nil); err != nil {
		log.Fatalf("Template execution error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
