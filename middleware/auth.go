package middleware

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strings"
)

type contextKey string

// RequestUserID is the authenticated user of the request
const RequestUserID contextKey = "request_user_id"

const (
	bearerPrefix = "Bearer "
	authHeader   = "Authorization"

	ErrTokenMissing = "Missing token in request headers"
	ErrTokenInvalid = "Request token is invalid"
)

type Auth struct {
	db *sql.DB
}

// NewAuth returns a new Auth with a db pool
func NewAuth(db *sql.DB) Auth {
	return Auth{db: db}
}

// ServeHTTP implements the negroni.Handler interface
func (a Auth) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	token := a.requestToken(r)

	if token == "" {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, "Error: ", ErrTokenMissing)
		// Exit middleware chain
		return
	}

	ctx := r.Context()
	userID, err := a.authorize(ctx, token)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "Error: ", ErrTokenMissing)
		default:
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "Error: ", err)
		}
		return
	}

	r = r.WithContext(context.WithValue(ctx, RequestUserID, userID))

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
