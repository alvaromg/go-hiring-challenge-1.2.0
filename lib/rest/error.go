package rest

import (
	"encoding/json"
	"net/http"

	"github.com/mytheresa/go-hiring-challenge/shared"
)

func HandleHTTPError(log shared.Logger, w http.ResponseWriter, r *http.Request, err error) {
	// log.WithContext(r.Context()).
	// 	WithFields(
	// 		logFieldsFromContext(r.Context()),
	// 	).
	// 	WithField("stack", liberr.GetErrStack(err)).
	// 	WithField("errorType", liberr.GetErrorType(err)).
	// 	Error(err.Error())

	w.Header().Add("Content-Type", "application/json")

	httpErrCode := errToHTTPCode(err)
	w.WriteHeader(httpErrCode)
	e := encodeErrorResponse(err)
	jsonOutput, err := json.Marshal(e)
	if err != nil {
		panic(err)
	}
	if _, err := w.Write(jsonOutput); err != nil {
		panic(err)
	}
}

func errToHTTPCode(err error) int {
	// if errors.Is(err, liberr.NewType(shared.ErrTypeAuthentication)) {
	// 	return http.StatusUnauthorized
	// }
	// if errors.Is(err, liberr.NewType(shared.ErrTypeAuthorization)) {
	// 	return http.StatusForbidden
	// }
	// if errors.Is(err, liberr.NewType(ErrTypeBadRequest)) ||
	// 	errors.Is(err, liberr.NewType(shared.ErrTypeInvalidApplicationInput)) ||
	// 	errors.Is(err, liberr.NewType(query.ErrTypeInvalidQuery)) {
	// 	return http.StatusBadRequest
	// }
	// if errors.Is(err, liberr.NewType(shared.ErrTypeNotFound)) {
	// 	return http.StatusNotFound
	// }
	// if errors.Is(err, liberr.NewType(shared.ErrTypeConflict)) ||
	// 	errors.Is(err, liberr.NewType(shared.ErrTypeDomainValidation)) {
	// 	return http.StatusConflict
	// }
	return http.StatusInternalServerError
}

type ErrorResponse struct {
	Metadata ResponseMetadata `json:"metadata"`
	Error    Error            `json:"error"`
}

type Error struct {
	Message string `json:"message"`
}

func encodeErrorResponse(err error) ErrorResponse {
	// errStr := ""
	// if errors.Is(err, liberr.NewType(shared.ErrTypeAuthentication)) {
	// 	errStr = shared.T(ctx, i18n.ErrorLibRestNotauthenticated)
	// } else {
	// 	errStr = err.Error()
	// }

	return ErrorResponse{
		Metadata: ResponseMetadata{
			// OperationId: shared.OperationIdFromContext(ctx),
		},
		Error: Error{
			Message: err.Error(),
		},
	}
}
