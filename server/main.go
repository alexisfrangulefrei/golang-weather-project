package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/efrei/weather/internal"
)

func main() {
	stations, err := internal.LoadFromJSON("sources/weather_data.json")

	if err != nil {
		log.Fatal(err)
	}

	store := NewStore()
	for _, s := range stations {
		store.Put(s)
	}
	log.Printf("bootstrap : %d stations chargées", len(stations))

	app := &App{store: store}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "ok")
	})
	mux.HandleFunc("GET /stations", app.listStations)
	mux.HandleFunc("GET /stations/{id}", app.getStation)
	mux.HandleFunc("POST /stations", app.createStation)
	mux.HandleFunc("PUT /stations/{id}", app.updateStation)

	log.Fatal(http.ListenAndServe(":8080", mux))
}
