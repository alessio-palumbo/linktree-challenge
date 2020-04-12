package links

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	e "github.com/alessio-palumbo/linktree-challenge/errors"
	"github.com/alessio-palumbo/linktree-challenge/handlers"
	"github.com/alessio-palumbo/linktree-challenge/handlers/models"
	"github.com/alessio-palumbo/linktree-challenge/middleware"
	"github.com/google/uuid"
)

const defaultOrder = "asc"

var validOrderKeys = map[string]bool{
	"created_at": true,
	"title":      true,
	"type":       true,
}

// IndexHandler list all the links for a given user.
type IndexHandler handlers.Group

func (h IndexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := middleware.CtxReqUserID(ctx)

	links, err := getUserLinks(ctx, h.DB, userID, r.FormValue("sort_by"))
	if err != nil {
		e.WriteError(w, http.StatusInternalServerError, err)
	}

	handlers.WriteResponse(w, http.StatusOK, links)
}

func getUserLinks(ctx context.Context, db *sql.DB, userID, sortBy string) ([]models.Link, error) {

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

	// TODO add default ordering based on a orderID, to be controlled by a different api/table
	if orderBy := sortByClause(sortBy); orderBy != "" {
		stmt += fmt.Sprintf(" ORDER BY %s ", orderBy)
	}

	rows, err := db.QueryContext(ctx, stmt, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	links := []models.Link{}
	for rows.Next() {
		var (
			l        models.Link
			subID    *uuid.UUID
			metadata *json.RawMessage
		)

		err := rows.Scan(&l.ID, &l.Type, &l.Title, &l.URL,
			&l.Thumbnail, &l.CreatedAt, &subID, &metadata)
		if err != nil {
			return nil, err
		}

		// Sublinks must have metadata as it is a required field
		if subID != nil && metadata != nil {
			_, err := addSublink(&l, (*subID).String(), *metadata)
			if err != nil {
				return nil, err
			}
		}

		links = append(links, l)
	}

	return links, rows.Err()
}

func sortByClause(sortBy string) string {
	if sortBy == "" {
		return sortBy
	}

	clauses := strings.Split(sortBy, ",")
	var sortClauses []string

	for _, c := range clauses {
		col := strings.Split(c, ":")
		if _, valid := validOrderKeys[col[0]]; valid {
			order := defaultOrder
			if len(col) > 1 && (col[1] == "asc" || col[1] == "desc") {
				order = col[1]
			}
			sortClauses = append(sortClauses, fmt.Sprintf("%s %s", col[0], order))
		}
	}

	return strings.Join(sortClauses, ", ")
}
