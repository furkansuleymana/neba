package handlers

import (
	"html/template"
	"net/http"

	"github.com/furkansuleymana/neba/network"
	"github.com/furkansuleymana/neba/ui"
)

func RegisterDevicesListRoute(mux *http.ServeMux) {
	mux.HandleFunc("/devices", handleDevicesList)
}

func handleDevicesList(w http.ResponseWriter, r *http.Request) {
	// Discover devices
	deviceList, err := network.DiscoverSSDP()
	if err != nil {
		http.Error(w, "Failed to discover devices: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Parse template from embedded filesystem
	tmpl, err := template.ParseFS(ui.TemplatesDirFS, "devices.html")
	if err != nil {
		http.Error(w, "Failed to parse template: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Execute template with device list data
	err = tmpl.Execute(w, deviceList)
	if err != nil {
		http.Error(w, "Failed to render template: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
