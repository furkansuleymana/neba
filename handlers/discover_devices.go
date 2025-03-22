package handlers

import (
	"html/template"
	"net/http"

	"github.com/furkansuleymana/neba/network"
	"github.com/furkansuleymana/neba/ui"
)

// DevicePageData contains the data for the discover_devices page
type DevicePageData struct {
	Title       string
	CurrentPath string
	NavItems    []NavItem
	Devices     []map[string]string
}

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

	// Prepare data for template
	data := DevicePageData{
		Title:       "Discover Devices",
		CurrentPath: r.URL.Path,
		NavItems:    navigationItems,
		Devices:     deviceList,
	}

	// Parse templates from embedded filesystem
	tmpl, err := template.ParseFS(ui.TemplatesDirFS,
		"layouts/base.html",
		"components/navigation.html",
		"routes/discover_devices.html")
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
