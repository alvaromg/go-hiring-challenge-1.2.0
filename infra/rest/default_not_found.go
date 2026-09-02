package rest

import (
	"encoding/json"
	"errors"
	"net/http"
)

func DefaultNotFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	w.WriteHeader(http.StatusNotFound)
	e := encodeErrorResponse(r.Context(), errors.New("route not found"))
	jsonOutput, err := json.Marshal(e)
	if err != nil {
		panic(err)
	}
	if _, err := w.Write(jsonOutput); err != nil {
		panic(err)
	}
}
