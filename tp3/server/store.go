package main

import (
	"maps"
	"slices"

	"github.com/efrei/weather/shared"
)

type Store struct {
	Stations map[string]shared.Station
}

func NewStore() *Store {
	return &Store{
		Stations: make(map[string]shared.Station),
	}
}

func (s *Store) Put(st shared.Station) {
	s.Stations[st.Id] = st
}

func (s *Store) Has(id string) bool {
	_, ok := s.Stations[id]
	return ok
}

func (s *Store) Get(id string) (shared.Station, bool) {
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

func (s *Store) All() []shared.Station {
	return slices.Collect(maps.Values(s.Stations))
}
