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
	// TODO: What if err == nil?
	apihErr, ok := err.(*Error)
	if !ok {
		apihErr = &Error{
			StatusCode: http.StatusInternalServerError,
			Message:    "Internal server error",
			Inner:      err,
		}
	} else if apihErr.Message == "" {
		if apihErr.Inner != nil {
			apihErr.Message = apihErr.Inner.Error()
		} else {
			apihErr.Message = http.StatusText(apihErr.StatusCode)
		}
	}
	ServeJSON(w, apihErr.StatusCode, apihErr)
}
