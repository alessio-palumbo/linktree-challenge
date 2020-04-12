package links

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	e "github.com/alessio-palumbo/linktree-challenge/errors"
	"github.com/alessio-palumbo/linktree-challenge/handlers"
	"github.com/alessio-palumbo/linktree-challenge/handlers/models"
	"github.com/alessio-palumbo/linktree-challenge/middleware"
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

	link, sublinks, err := prepareDbObject(body, h.Validator)
	if err != nil {
		e.WriteError(w, http.StatusBadRequest, err)
		return
	}

	err = h.insertLinks(r.Context(), link, sublinks)
	if err != nil {
		e.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	handlers.WriteResponse(w, http.StatusCreated, *link)
}

func prepareDbObject(body []byte, validator *validator.CustomValidator) (*models.Link, []models.Sublink, error) {

	var l models.LinkPayload
	err := json.Unmarshal(body, &l)
	if err := e.CheckValid(err, l, validator); err != nil {
		return nil, nil, err
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
			subID, ID := models.GenerateUUIDPair()

			sl, err := addSublink(link, ID, s)
			if err := e.CheckValid(err, sl, validator); err != nil {
				return nil, nil, err
			}

			data, err := json.Marshal(sl)
			if err != nil {
				return nil, nil, err
			}

			dbSub := models.Sublink{ID: subID, Metadata: data}
			dbSubs = append(dbSubs, dbSub)
		}

		return link, dbSubs, nil
	}

	return link, nil, nil
}

func (h *PostHandler) insertLinks(ctx context.Context, l *models.Link, sl []models.Sublink) error {
	userID := middleware.CtxReqUserID(ctx)

	tx, err := h.DB.Begin()
	if err != nil {
		return err
	}

	l.UUID, l.ID = models.GenerateUUIDPair()

	_, err = tx.ExecContext(ctx, `
		INSERT INTO links (id, user_id, type, title, url, thumbnail)
		VALUES ($1, $2, $3, $4, $5, $6)
		`, l.UUID, userID, l.Type, l.Title, l.URL, l.Thumbnail)

	if err != nil {
		tx.Rollback()
		return err
	}

	if len(sl) > 0 {
		stmt, values := generateBulkInsert(sl)

		_, err = tx.ExecContext(ctx, stmt, values...)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func generateBulkInsert(sl []models.Sublink) (string, []interface{}) {

	cols := 3
	values := make([]interface{}, 0, len(sl)*cols)
	stmt := strings.TrimRight(fmt.Sprintf(`
		INSERT INTO sublinks (id, link_id, metadata) VALUES %s`,
		strings.Repeat(" ($1, $2, $3),", len(sl))), ",")

	for _, s := range sl {
		values = append(values, s.ID, s.UserID, s.Metadata)
	}

	return stmt, values
}
