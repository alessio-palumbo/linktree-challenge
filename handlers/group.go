package handlers

import (
	"database/sql"

	"github.com/alessio-palumbo/linktree-challenge/middleware"
)

// Group injects DB pool and Auth client in handler
type Group struct {
	Auth middleware.Auth
	DB   *sql.DB
}
