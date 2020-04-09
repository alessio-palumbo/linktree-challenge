package server

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx"
)

// New returns a handler to serve the links api.
func New(db *pgx.ConnPool) http.Handler {
	// TODO Add middlewares

	// Add mux and healthcheck
	router := mux.NewRouter()
	linksSB := router.
		PathPrefix("/api/links").
		Subrouter()

	linksSB.Handle("", nil).Methods("GET")
	linksSB.Handle("/{link_id}", nil).Methods("GET")
	linksSB.Handle("", nil).Methods("POST")
	linksSB.Handle("/{link_id}", nil).Methods("PUT")
	linksSB.Handle("/{link_id}", nil).Methods("DELETE")

	sublinksSB := router.
		PathPrefix("/api/sublinks").
		Subrouter()

	sublinksSB.Handle("", nil).Methods("POST")
	sublinksSB.Handle("/{link_id}", nil).Methods("PUT")
	sublinksSB.Handle("/{link_id}", nil).Methods("DELETE")

	return router
}
