package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/furkansuleymana/neba/network"
)

func RegisterFindDevicesRoute(mux *http.ServeMux) {
	mux.HandleFunc("/api/find_devices", handleFindDevices)
}

func handleFindDevices(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	devices, err := network.DiscoverSSDP()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(devices); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
