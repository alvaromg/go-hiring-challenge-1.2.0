package rest

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/mytheresa/go-hiring-challenge/shared"

	domainerrors "github.com/mytheresa/go-hiring-challenge/domain/errors"
	"github.com/mytheresa/go-hiring-challenge/domain/query"

	libmodel "github.com/mytheresa/go-hiring-challenge/infra/database"
	"github.com/mytheresa/go-hiring-challenge/infra/operation"
)

var (
	ErrorBadRequest = errors.New("bad request error")
)

func HandleHTTPError(log shared.Logger, w http.ResponseWriter, r *http.Request, err error) {

	log.WithContext(r.Context()).Errorf("%s", err.Error())
	if errors.Is(err, libmodel.ErrorPersistence) {
		// if it's a persistence error don't include details in http response
		err = libmodel.ErrorPersistence
	}

	w.Header().Add("Content-Type", "application/json")

	httpErrCode := errToHTTPCode(err)
	w.WriteHeader(httpErrCode)
	e := encodeErrorResponse(r.Context(), err)
	jsonOutput, err := json.Marshal(e)
	if err != nil {
		panic(err)
	}
	if _, err := w.Write(jsonOutput); err != nil {
		panic(err)
	}
}

func errToHTTPCode(err error) int {
	if errors.Is(err, ErrorBadRequest) ||
		errors.Is(err, query.ErrorInvalidQuery) ||
		errors.Is(err, domainerrors.ErrorDomainValidation) {
		return http.StatusBadRequest
	}
	if errors.Is(err, domainerrors.ErrorNotFound) {
		return http.StatusNotFound
	}
	if errors.Is(err, domainerrors.ErrorDuplicatedResource) {
		return http.StatusConflict
	}
	return http.StatusInternalServerError
}

type ErrorResponse struct {
	Metadata ResponseMetadata `json:"metadata"`
	Error    Error            `json:"error"`
}

type Error struct {
	Message string `json:"message"`
}

func encodeErrorResponse(ctx context.Context, err error) ErrorResponse {
	return ErrorResponse{
		Metadata: ResponseMetadata{
			OperationId: operation.IdFromContext(ctx),
		},
		Error: Error{
			Message: err.Error(),
		},
	}
}
