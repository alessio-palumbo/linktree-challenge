package handlers

import (
	"database/sql"
	"net/http"
)

// Group injects DB pool and Auth client in handler
type Group struct {
	Auth http.Client
	DB   *sql.DB
}
