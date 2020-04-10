package middleware

import (
	"context"
	"database/sql"
	"net/http"
	"strings"

	e "github.com/alessio-palumbo/linktree-challenge/errors"
)

const (
	bearerPrefix = "Bearer "
	authHeader   = "Authorization"

	errTokenMissing = "missing token in request headers"
	errTokenInvalid = "request token is invalid"
)

type Auth struct {
	db *sql.DB
}

// NewAuth returns a new Auth with a db pool
func NewAuth(db *sql.DB) Auth {
	return Auth{db: db}
}

// ServeHTTP implements the negroni.Handler interface
// TODO this could be using a jwt token containing the userID and any other data needed,
// while using the db tokens table to check the expiry timestamp.
func (a Auth) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	token := a.requestToken(r)

	if token == "" {
		e.WriteError(w, http.StatusUnauthorized, errTokenMissing)
		return
	}

	ctx := r.Context()
	userID, err := a.authorize(ctx, token)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			e.WriteError(w, http.StatusUnauthorized, errTokenInvalid)
		default:
			e.WriteError(w, http.StatusInternalServerError, err)
		}
		return
	}

	r = CtxSetUserID(ctx, r, userID)

	next(w, r)
}

func (a Auth) requestToken(r *http.Request) string {
	var token string
	auth := r.Header.Get(authHeader)

	if strings.HasPrefix(auth, bearerPrefix) {
		token = auth[len(bearerPrefix):]
	}

	return token
}

func (a Auth) authorize(ctx context.Context, token string) (string, error) {
	stmt := `
		SELECT user_id
		  FROM user_tokens
		 WHERE id = $1
		   AND expire_at > NOW()
	`

	var userID string
	err := a.db.QueryRowContext(ctx, stmt, token).Scan(&userID)
	return userID, err
}
