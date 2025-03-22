package handlers

import (
	"html/template"
	"net/http"

	"github.com/furkansuleymana/neba/network"
	"github.com/furkansuleymana/neba/ui"
)

func RegisterDiscoverDevicesRoute(mux *http.ServeMux) {
	mux.HandleFunc("/discover_devices", handleDiscoverDevices)
}

func handleDiscoverDevices(w http.ResponseWriter, r *http.Request) {
	// Discover devices
	deviceList, err := network.DiscoverSSDP()
	if err != nil {
		http.Error(w, "Failed to discover devices: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Parse template from embedded filesystem
	template, err := template.ParseFS(ui.TemplatesDirFS, "discover_devices.html")
	if err != nil {
		http.Error(w, "Failed to parse template: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Execute template with device list data
	err = template.Execute(w, deviceList)
	if err != nil {
		http.Error(w, "Failed to render template: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
