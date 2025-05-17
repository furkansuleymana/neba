package handlers

import (
	"html/template"
	"log"
	"net/http"

	"github.com/furkansuleymana/neba/ui"
)

var (
	manageDevicesTmpl *template.Template
)

func RegisterManageDevicesRoute(fs http.Handler, mux *http.ServeMux) {
	var err error
	manageDevicesTmpl, err = template.ParseFS(ui.FS, "manage.html")
	if err != nil {
		log.Fatalf("Failed to parse templates: %v", err)
	}

	mux.HandleFunc("/manage", handleManageDevices)
}

func handleManageDevices(w http.ResponseWriter, r *http.Request) {
	if err := manageDevicesTmpl.ExecuteTemplate(w, "manage.html", nil); err != nil {
		log.Fatalf("Template execution error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
