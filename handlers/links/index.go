package links

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"

	e "github.com/alessio-palumbo/linktree-challenge/errors"
	"github.com/alessio-palumbo/linktree-challenge/handlers"
	"github.com/alessio-palumbo/linktree-challenge/handlers/models"
	"github.com/alessio-palumbo/linktree-challenge/middleware"
)

// IndexHandler list all the links for a given user.
type IndexHandler handlers.Group

func (h IndexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := middleware.CtxReqUserID(ctx)

	links, err := getUserLinks(ctx, h.DB, userID)
	if err != nil {
		e.WriteError(w, http.StatusInternalServerError, err)
	}

	json.NewEncoder(w).Encode(links)
}

func getUserLinks(ctx context.Context, db *sql.DB, userID string) ([]models.Link, error) {

	stmt := `
		SELECT l.id,
		       l.type,
		       l.title,
		       l.url,
		       l.thumbnail,
		       l.created_at,

		       sl.id,
		       sl.metadata
		  FROM links l
		  LEFT JOIN sublinks sl ON sl.link_id = l.id
		 WHERE l.user_id = $1
	`

	rows, err := db.QueryContext(ctx, stmt, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	links := []models.Link{}

	for rows.Next() {
		var l models.Link
		var sl models.Sublink

		err := rows.Scan(&l.ID, &l.Type, &l.Title, &l.URL,
			&l.CreatedAt, &sl.ID, &sl.Metadata)
		if err != nil {
			return nil, err
		}

		switch l.Type {
		case models.LinkMusic:
			// TODO
		case models.LinkShows:
			// TODO
		}

		links = append(links, l)
	}

	return links, rows.Err()
}
