package links

import (
	"net/http"

	"github.com/alessio-palumbo/linktree-challenge/handlers"
)

// IndexHandler list all the links for a given user.
type IndexHandler handlers.Group

func (h IndexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// TODO
}
