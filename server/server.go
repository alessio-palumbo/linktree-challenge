package server

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"

	"github.com/alessio-palumbo/linktree-challenge/handlers"
)

// New returns a handler to serve the links api.
func New(g handlers.Group) http.Handler {
	// Add middleware classic package with recover and logging
	n := negroni.Classic()

	// // Add endpoint to check db connection
	n.UseFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		if r.URL.Path == "/healthcheck" {
			var ok bool
			if err := g.DB.QueryRowContext(r.Context(), "SELECT true").Scan(&ok); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprint(w, "Error: ", err)
			}
			fmt.Fprint(w, "OK")
			return
		}

		next(w, r)
	})

	// TODO Register auth middleware

	// Add multiplexer and register routes
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

	n.UseHandler(router)

	return n
}
