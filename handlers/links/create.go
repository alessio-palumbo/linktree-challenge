package links

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	e "github.com/alessio-palumbo/linktree-challenge/errors"
	"github.com/alessio-palumbo/linktree-challenge/handlers"
	"github.com/alessio-palumbo/linktree-challenge/handlers/models"
	"github.com/alessio-palumbo/linktree-challenge/validator"
)

// PostHandler list all the links for a given user.
type PostHandler handlers.Group

func (h PostHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		e.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	link, err := prepareDbObject(body, h.Validator)
	if err != nil {
		e.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// err = h.insertLinks(r.Context(), link)
	// if err != nil {
	// 	e.WriteError(w, http.StatusInternalServerError, err)
	// }

	json.NewEncoder(w).Encode(*link)
}

func prepareDbObject(body []byte, validator *validator.CustomValidator) (*models.Link, error) {

	var l models.LinkPayload
	err := json.Unmarshal(body, &l)
	if err := e.CheckValid(err, l, validator); err != nil {
		return nil, err
	}

	link := &models.Link{
		Type:      l.Type,
		Title:     l.Title,
		Thumbnail: l.Thumbnail,
		URL:       l.URL,
	}

	if len(l.SubLinks) > 0 {
		dbSubs := make([]models.Sublink, 0, len(l.SubLinks))

		for _, s := range l.SubLinks {
			//dbSub.ID, ID = models.GenerateUUIDPair()
			sl, err := addSublink(link, "", s)
			if err := e.CheckValid(err, sl, validator); err != nil {
				return nil, err
			}

			data, err := json.Marshal(sl)
			if err != nil {
				return nil, err
			}

			dbSub := models.Sublink{Metadata: data}
			dbSubs = append(dbSubs, dbSub)
		}
	}

	return link, nil
}
