package errors

import (
	"fmt"
	"net/http"

	"github.com/alessio-palumbo/linktree-challenge/validator"
)

// TODO define an error response struct and methods

// WriteError prints an json error with the given status and message
func WriteError(w http.ResponseWriter, status int, err interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	fmt.Fprint(w, JSONError(err))
}

// JSONError formats an error to a json response
func JSONError(err interface{}) string {
	return fmt.Sprintf(`{"error":%q}`, err)
}

// CheckValid is a helper that checks the error coming from an unmarshalling and
// validate the interface through the custom validator
func CheckValid(err error, i interface{}, cv *validator.CustomValidator) error {
	if err != nil {
		return err
	}

	err = cv.Validate(i)
	return err
}
