package handlers

import (
	"net/http"

	"github.com/furkansuleymana/neba/api"
)

func RegisterParamRoute(mux *http.ServeMux) {
	mux.HandleFunc("/api/param", paramHandler())
}

func paramHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		err := api.Param("ip_address", "username", "password")

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
