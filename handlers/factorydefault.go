package handlers

import (
	"net/http"

	"github.com/furkansuleymana/neba/api"
)

func RegisterFactoryDefaultRoute(mux *http.ServeMux) {
	mux.HandleFunc("/api/factorydefault", factoryDefaultHandler())
}

func factoryDefaultHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		err := api.FactoryDefault()

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
