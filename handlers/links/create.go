package links

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	e "github.com/alessio-palumbo/linktree-challenge/errors"
	"github.com/alessio-palumbo/linktree-challenge/handlers"
	"github.com/alessio-palumbo/linktree-challenge/handlers/models"
)

// PostHandler list all the links for a given user.
type PostHandler handlers.Group

func (h PostHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		e.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	var l models.Link
	err = json.Unmarshal(body, &l)
	if err != nil {
		e.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	err = h.Validator.Validate(l)
	if err != nil {
		e.WriteError(w, http.StatusBadRequest, err)
		return
	}

	json.NewEncoder(w).Encode(l)
}
