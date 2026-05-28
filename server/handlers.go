package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/efrei/weather/internal"
)

type App struct{ store *Store }

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(v)

	if err != nil {
		log.Fatal(err)
	}
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func (a *App) listStations(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(a.store.All())

	if err != nil {
		log.Fatal(err)
	}
}

func (a *App) getStation(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	st, ok := a.store.Get(id)

	if !ok {
		writeError(w, http.StatusNotFound,
			fmt.Sprintf("station %q introuvable", id))
		return
	}

	writeJSON(w, http.StatusOK, st)
}

func (a *App) createStation(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var st internal.Station

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(&st); err != nil {
		writeError(w, 400, "JSON invalide: "+err.Error())
		return
	}

	if st.Id == "" {
		writeError(w, 400, "id manquant")
		return
	}

	if a.store.Has(st.Id) {
		writeError(w, 409, "id "+st.Id+" déjà utilisé")
		return
	}

	a.store.Put(st)

	writeJSON(w, 201, st)
}
