package handlers

import (
	"html/template"
	"log"
	"net/http"

	"github.com/furkansuleymana/neba/network"
	"github.com/furkansuleymana/neba/ui"
)

var (
	findDevicesTmpl *template.Template
)

// DevicePageData contains the data for the find_devices page
type DevicePageData struct {
	Devices     []map[string]string
	DeviceCount int
	Error       string
}

func RegisterFindDevicesRoute(fs http.Handler, mux *http.ServeMux) {
	var err error
	findDevicesTmpl, err = template.ParseFS(ui.FS, "find_devices.html")
	if err != nil {
		log.Fatalf("Failed to parse templates: %v", err)
	}

	mux.HandleFunc("/find_devices", handleFindDevices)
}

func handleFindDevices(w http.ResponseWriter, r *http.Request) {
	data := DevicePageData{}

	deviceList, err := network.FindSSDP()
	if err != nil {
		data.Error = err.Error()
	}

	data.Devices = deviceList
	data.DeviceCount = len(deviceList)

	if err := findDevicesTmpl.ExecuteTemplate(w, "find_devices.html", data); err != nil {
		log.Fatalf("Template execution error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
