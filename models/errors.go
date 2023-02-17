package models

import (
	"errors"
	"net/http"
)

type RequestError struct {
	StatusCode int
	Err        error
}

var ErrDocumentNotFound = errors.New("DocumentNotFound")

func NewErrorStatusCodeMaps() map[error]int {

	var errorStatusCodeMaps = make(map[error]int)
	errorStatusCodeMaps[ErrDocumentNotFound] = http.StatusNotFound
	return errorStatusCodeMaps
}
