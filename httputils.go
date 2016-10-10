package main

import (
	"crypto/subtle"
	"fmt"
	"log"
	"mime"
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

func parseRequestContentType(r *http.Request) (string, error) {
	ctype := r.Header.Get("Content-Type")
	if ctype == "" {
		return "", nil
	}

	mediatype, _, err := mime.ParseMediaType(ctype)
	if err != nil {
		return "", newHTTPError(http.StatusBadRequest, "Failed to parse Content-Type: %v", err)
	}
	return mediatype, nil
}

func verifyToken(token string, correctToken string) error {
	if token == "" {
		return newHTTPError(http.StatusUnauthorized, "Missing token")
	}
	if subtle.ConstantTimeCompare([]byte(token), []byte(correctToken)) != 1 {
		return newHTTPError(http.StatusUnauthorized, "Invalid token")
	}
	return nil
}
