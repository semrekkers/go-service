package apih

import (
	"fmt"
	"net/http"
)

type Error struct {
	StatusCode int    `json:"-"`
	Message    string `json:"message"`
	Inner      error  `json:"-"`
}

func (e *Error) Error() string {
	if e.Inner != nil {
		return fmt.Sprintf("%s: %s", e.Message, e.Inner.Error())
	}
	return e.Message
}

func (e *Error) Unwrap() error {
	return e.Inner
}

func ServeError(w http.ResponseWriter, err error) {
	apihErr, ok := err.(*Error)
	if !ok {
		apihErr = &Error{
			StatusCode: http.StatusInternalServerError,
			Message:    "Internal server error",
			Inner:      err,
		}
	}
	ServeJSON(w, apihErr.StatusCode, apihErr)
}
