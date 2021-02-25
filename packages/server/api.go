package server

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jorgemhc-dev/academy-go-q12021/packages/getinfo"
)

type api struct {
	router http.Handler
}

type Server interface {
	Router() http.Handler
}


func (a *api) Router() http.Handler {
	return a.router
}

func New() Server {
	a := &api{}

	r := mux.NewRouter()

	r.HandleFunc("/pokemon", getinfo.FetchAll).Methods(http.MethodGet)
	r.HandleFunc("/pokemon/{ID:[A-Za-z0-9]+}", getinfo.FetchPokemon).Methods(http.MethodGet)

	a.router = r
	return a
}
