package server

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"

	e "github.com/alessio-palumbo/linktree-challenge/errors"
	"github.com/alessio-palumbo/linktree-challenge/handlers"
	"github.com/alessio-palumbo/linktree-challenge/handlers/links"
)

// New returns a handler to serve the links api.
func New(g handlers.Group) http.Handler {
	// Add middleware classic package with recover and logging
	n := negroni.Classic()

	// Add endpoint to check db connection
	n.UseFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		if r.URL.Path == "/healthcheck" {
			var ok bool
			err := g.DB.QueryRowContext(r.Context(), "SELECT true as ok").Scan(&ok)
			switch err {
			case nil:
				fmt.Fprint(w, "OK")
			default:
				e.WriteError(w, http.StatusInternalServerError, err)
			}
			return
		}

		next(w, r)
	})

	// Add authentication to middleware chain
	n.Use(g.Auth)

	// Add multiplexer and register routes
	router := mux.NewRouter()

	linksSB := router.
		PathPrefix("/api/links").
		Subrouter()

	linksSB.Handle("", links.IndexHandler(g)).Methods("GET")
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
