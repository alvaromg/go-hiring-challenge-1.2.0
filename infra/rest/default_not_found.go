package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mytheresa/go-hiring-challenge/shared"
)

func DefaultNotFound(log shared.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")

		w.WriteHeader(http.StatusNotFound)

		path := r.URL.Path

		e := encodeErrorResponse(r.Context(), fmt.Errorf("route %q not found", path))
		jsonOutput, err := json.Marshal(e)
		if err != nil {
			log.WithContext(r.Context()).Errorf("error encoding http error response: %s", err)
		}
		if _, err := w.Write(jsonOutput); err != nil {
			log.WithContext(r.Context()).Errorf("error writing http error response: %s", err)
		}
	}

}
