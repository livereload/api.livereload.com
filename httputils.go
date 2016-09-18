package main

import (
	"crypto/subtle"
	"fmt"
	"log"
	"net/http"
)

type httpError struct {
	code    int
	message string
}

func (e httpError) Error() string {
	return fmt.Sprintf("%s (HTTP %d)", e.message, e.code)
}

func newHTTPError(code int, format string, a ...interface{}) error {
	message := fmt.Sprintf(format, a...)
	return &httpError{code, message}
}

func sendErrorMessage(w http.ResponseWriter, code int, message string) {
	if message == "" {
		message = http.StatusText(code)
	}
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(code)
	w.Write([]byte(message))
}

func sendErrorFmt(w http.ResponseWriter, code int, format string, a ...interface{}) {
	sendErrorMessage(w, code, fmt.Sprintf(format, a...))
}

func sendError(w http.ResponseWriter, err error) {
	if httperr, ok := err.(*httpError); ok {
		sendErrorMessage(w, httperr.code, httperr.message)
	} else {
		log.Printf("** Unexpected error: %v", err)
		sendErrorMessage(w, http.StatusInternalServerError, "")
	}
}

func verifyToken(r *http.Request, correctToken string) error {
	token := r.Header.Get("X-Token")
	if token == "" {
		return newHTTPError(http.StatusUnauthorized, "Missing auth token")
	}
	if subtle.ConstantTimeCompare([]byte(token), []byte(correctToken)) != 1 {
		return newHTTPError(http.StatusUnauthorized, "Invalid auth token")
	}
	return nil
}
