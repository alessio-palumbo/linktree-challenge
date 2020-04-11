package handlers

import (
	"database/sql"

	"github.com/alessio-palumbo/linktree-challenge/middleware"
	"github.com/alessio-palumbo/linktree-challenge/validator"
)

// Group injects DB pool and Auth client in handler
type Group struct {
	Auth      middleware.Auth
	DB        *sql.DB
	Validator *validator.CustomValidator
}
