package errors

import (
	"fmt"
	"net/http"
)

// TODO define an error response struct and methods

// WriteError prints an json error with the given status and message
func WriteError(w http.ResponseWriter, status int, err interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	fmt.Fprintf(w, `{"error":%q}`, err)
}
