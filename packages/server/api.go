package server

import (
	"net/http"

	"github.com/gorilla/mux"
)

type api struct {
	router http.Handler
}

type Server interface {
	Router() http.Handler
}

func New() Server {
	a := &api{}

	r := mux.NewRouter()

	r.HandleFunc("/pokemons", a.fetchAll).Methods(http.MethodGet)
	r.HandleFunc("/pokemon/{ID:[0-9]+}", a.fetchPokemon).Methods(http.MethodGet)

	a.router = r
	return a
}

func (a *api) Router() http.Handler {
	return a.router
}

func (a *api) fetchAll(w http.ResponseWriter, r *http.Request) {}
func (a *api) fetchPokemon(w http.ResponseWriter, r *http.Request) {}