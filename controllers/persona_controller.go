package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/danysoftdev/p-go-search/services"

	"github.com/gorilla/mux"
)

func ObtenerPersonaPorDocumento(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	documento := params["documento"]

	persona, err := services.BuscarPersonaPorDocumento(documento)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(persona)
}
