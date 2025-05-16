package handlers

import (
	"html/template"
	"log"
	"net/http"

	"github.com/furkansuleymana/neba/network"
	"github.com/furkansuleymana/neba/ui"
)

// DevicePageData contains the data for the find_devices page
type DevicePageData struct {
	Devices []map[string]string
}

func RegisterFindDevicesRoute(mux *http.ServeMux) {
	mux.HandleFunc("/find_devices", handleFindDevices)
}

func handleFindDevices(w http.ResponseWriter, r *http.Request) {
	// Discover devices
	deviceList, err := network.FindSSDP()
	if err != nil {
		log.Println("Error discovering devices:", err)
	}

	// Prepare data for template
	data := DevicePageData{
		Devices: deviceList,
	}

	// Parse templates from embedded filesystem
	tmpl, err := template.ParseFS(ui.TemplatesDirFS,
		"index.html",
		"find_devices.html")
	if err != nil {
		http.Error(w, "Failed to parse templates: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Execute template with device list data
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Failed to render template: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
