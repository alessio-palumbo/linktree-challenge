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

	var l models.LinkPayload
	err = json.Unmarshal(body, &l)
	if err := e.CheckValid(err, l, h.Validator); err != nil {
		e.WriteError(w, http.StatusBadRequest, err)
		return
	}

	link := models.Link{
		Type:      l.Type,
		Title:     l.Title,
		Thumbnail: l.Thumbnail,
		URL:       l.URL,
	}

	if len(l.SubLinks) > 0 {
		sublinks, err := subLinks(&l, h.Validator)
		if err != nil {
			e.WriteError(w, http.StatusBadRequest, err)
			return
		}

		for _, s := range sublinks {
			addSublink(&link, "", &s.Metadata)
		}
	}

	// links, err := getUserLinks(ctx, h.DB, userID, r.FormValue("sort_by"))
	// if err != nil {
	// 	e.WriteError(w, http.StatusInternalServerError, err)
	// }

	json.NewEncoder(w).Encode(link)
}

// subLinks validates data by unmarshalling it into the correct model and run
// validator before building the json that will form the DB object
// Note: If the sublink payload matches any, but not all the fields of the model, the matching fields
// 	 will still be parsed and the sublink considered correct, as long as it satisfy field validation
// TODO this could be used by a PUT call so change it to exported type in a shared folder
func subLinks(l *models.LinkPayload, validator *validator.CustomValidator) ([]models.Sublink, error) {
	sublinks := make([]models.Sublink, 0, len(l.SubLinks))
	var err error

	for _, sub := range l.SubLinks {
		var subDB models.Sublink
		var data json.RawMessage

		switch l.Type {
		case models.LinkMusic:
			sl := models.Platform{}
			err = json.Unmarshal(sub, &sl)
			if err := e.CheckValid(err, sl, validator); err != nil {
				return nil, err
			}

			data, err = json.Marshal(sl)
		case models.LinkShows:
			sl := models.Show{}
			err = json.Unmarshal(sub, &sl)
			if err := e.CheckValid(err, sl, validator); err != nil {
				return nil, err
			}

			data, err = json.Marshal(sl)
		}

		if err != nil {
			return nil, err
		}

		subDB.Metadata = data
		sublinks = append(sublinks, subDB)

	}

	return sublinks, nil
}
