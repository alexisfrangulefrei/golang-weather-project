package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/efrei/weather/shared"
)

type App struct{ store *Store }

type ErrorResponse struct {
	Error string `json:"error"`
	Code  string `json:"code,omitempty"`
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(v)

	if err != nil {
		log.Fatal(err)
	}
}

func writeError(w http.ResponseWriter, status int, code string, msg string) {
	writeJSON(w, status, ErrorResponse{Error: msg, Code: code})
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
		writeError(w, http.StatusNotFound, "NOT_FOUND", fmt.Sprintf("station %q introuvable", id))
		return
	}

	writeJSON(w, http.StatusOK, st)
}

func (a *App) createStation(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var st shared.Station

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(&st); err != nil {
		writeError(w, http.StatusBadRequest, "BAD_JSON", "JSON invalide: "+err.Error())
		return
	}

	if st.Id == "" {
		writeError(w, http.StatusBadRequest, "BAD_JSON", "id manquant")
		return
	}

	if a.store.Has(st.Id) {
		writeError(w, http.StatusConflict, "ID_TAKEN", "id "+st.Id+" déjà utilisé")
		return
	}

	a.store.Put(st)

	writeJSON(w, http.StatusCreated, st)
}

func (a *App) updateStation(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	id := r.PathValue("id")
	var st shared.Station

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(&st); err != nil {
		writeError(w, http.StatusBadRequest, "BAD_JSON", "JSON invalide: "+err.Error())
		return
	}

	if st.Id != "" && st.Id != id {
		writeError(w, http.StatusBadRequest, "INCOHERENCE_BODY_URL_ID", "incohérence id body vs URL")
		return
	}

	st.Id = id
	created := !a.store.Has(id)
	a.store.Put(st)
	if created {
		writeJSON(w, http.StatusCreated, st)
		return
	}

	writeJSON(w, http.StatusOK, st)
}

func (a *App) deleteStation(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if !a.store.Delete(id) {
		writeError(w, http.StatusNotFound, "NOT_FOUND", fmt.Sprintf("station %q introuvable", id))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (a *App) listObservations(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	st, ok := a.store.Get(id)
	if !ok {
		writeError(w, http.StatusNotFound, "NOT_FOUND", "station introuvable")
		return
	}
	writeJSON(w, http.StatusOK, st.Observations)
}
