package handlers

import (
	"html/template"
	"log"
	"net/http"

	"github.com/furkansuleymana/neba/network"
	"github.com/furkansuleymana/neba/ui"
)

var (
	discoverDevicesTmpl *template.Template
)

// DevicePageData contains the data for the /discover page
type DevicePageData struct {
	Devices     []map[string]string
	DeviceCount int
	Error       string
}

func RegisterDiscoverDevicesRoute(fs http.Handler, mux *http.ServeMux) {
	var err error
	discoverDevicesTmpl, err = template.ParseFS(ui.FS, "discover.html")
	if err != nil {
		log.Fatalf("Failed to parse templates: %v", err)
	}

	mux.HandleFunc("/discover", handleDiscoverDevices)
}

func handleDiscoverDevices(w http.ResponseWriter, r *http.Request) {
	data := DevicePageData{}

	deviceList, err := network.DiscoverSSDP()
	if err != nil {
		data.Error = err.Error()
	}

	data.Devices = deviceList
	data.DeviceCount = len(deviceList)

	if err := discoverDevicesTmpl.ExecuteTemplate(w, "discover.html", data); err != nil {
		log.Fatalf("Template execution error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
