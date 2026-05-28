package main

import (
	"maps"
	"slices"

	"github.com/efrei/weather/internal"
)

type Store struct {
	Stations map[string]internal.Station
}

func NewStore() *Store {
	return &Store{
		Stations: make(map[string]internal.Station),
	}
}

func (s *Store) Put(st internal.Station) {
	s.Stations[st.Id] = st
}

func (s *Store) Has(id string) bool {
	_, ok := s.Stations[id]
	return ok
}

func (s *Store) Get(id string) (internal.Station, bool) {
	st, ok := s.Stations[id]
	return st, ok
}

func (s *Store) Delete(id string) bool {
	if _, ok := s.Stations[id]; !ok {
		return false
	}
	delete(s.Stations, id)
	return true
}

func (s *Store) All() []internal.Station {
	return slices.Collect(maps.Values(s.Stations))
}
