package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/efrei/weather/internal"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "ok")
	})
	http.ListenAndServe(":8080", mux)

	stations, err := internal.LoadFromJSON("sources/weather_data.json")
	if err != nil {
		log.Fatal(err)
	}
	store := NewStore()
	for _, s := range stations {
		store.Put(s)
	}
	log.Printf("bootstrap : %d stations chargées", len(stations))
}
