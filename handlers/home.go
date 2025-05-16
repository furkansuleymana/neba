package handlers

import (
	"html/template"
	"log"
	"net/http"

	"github.com/furkansuleymana/neba/ui"
)

var (
	homeTmpl *template.Template
)

func RegisterHomeRoute(fs http.Handler, mux *http.ServeMux) {
	var err error
	homeTmpl, err = template.ParseFS(ui.FS, "index.html", "home.html")
	if err != nil {
		log.Fatalf("Failed to parse templates: %v", err)
	}

	mux.HandleFunc("/", handleHome)
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	if err := homeTmpl.ExecuteTemplate(w, "index.html", nil); err != nil {
		log.Printf("Template execution error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
