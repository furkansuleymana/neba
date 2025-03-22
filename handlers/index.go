package handlers

import (
	"html/template"
	"net/http"
	"strings"

	"github.com/furkansuleymana/neba/ui"
)

// PageInfo represents information about a page to be displayed in the index
type PageInfo struct {
	Name string
	Path string
}

func RegisterIndexRoute(mux *http.ServeMux) {
	mux.HandleFunc("/", handleIndex)
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	// If path is not root, let the file server handle it
	if r.URL.Path != "/" {
		http.FileServer(http.FS(ui.TemplatesDirFS)).ServeHTTP(w, r)
		return
	}

	// Read all HTML files from the embedded filesystem
	pages := []PageInfo{}

	entries, err := ui.TemplatesDirFS.ReadDir(".")
	if err != nil {
		http.Error(w, "Failed to read templates directory: "+err.Error(), http.StatusInternalServerError)
		return
	}

	for _, entry := range entries {
		// Skip directories and non-HTML files
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".html") || entry.Name() == "index.html" {
			continue
		}

		// Add page info to the list
		pages = append(pages, PageInfo{
			Name: strings.TrimSuffix(entry.Name(), ".html"),
			Path: "/" + entry.Name(),
		})
	}

	// Parse template from embedded filesystem
	tmpl, err := template.ParseFS(ui.TemplatesDirFS, "index.html")
	if err != nil {
		http.Error(w, "Failed to parse template: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Execute template with page list data
	err = tmpl.Execute(w, pages)
	if err != nil {
		http.Error(w, "Failed to render template: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
