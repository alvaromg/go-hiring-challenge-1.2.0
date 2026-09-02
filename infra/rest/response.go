package rest

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/mytheresa/go-hiring-challenge/shared"

	"github.com/mytheresa/go-hiring-challenge/infra/operation"
)

type JsonString json.RawMessage

type Response[T any] struct {
	Metadata ResponseMetadata `json:"metadata"`
	Data     T                `json:"data,omitempty"`
}

type ResponseMetadata struct {
	OperationId string            `json:"operationId"`
	Websocket   map[string]string `json:"websocket,omitempty"`
}

type EmptyResponse struct {
}

func encodeHandlerResponse[DO, RO any](log shared.Logger, ctx context.Context, encoder func(DO) (RO, error), out DO, w http.ResponseWriter, r *http.Request) {
	encodedData, err := encoder(out)
	if err != nil {
		HandleHTTPError(log, w, r, err)
		return
	}

	restOut := Response[RO]{
		Metadata: ResponseMetadata{
			OperationId: operation.IdFromContext(ctx),
		},
		Data: encodedData,
	}

	jsonOutput, err := json.Marshal(restOut)
	if err != nil {
		HandleHTTPError(log, w, r, err)
		return
	}
	if _, err := w.Write(jsonOutput); err != nil {
		HandleHTTPError(log, w, r, err)
	}
}

func EncodeResponse[T any](data T) Response[T] {
	return Response[T]{
		Data: data,
	}
}

func EncodeEmptyResponse(ctx context.Context, _ struct{}) (struct{}, error) {
	return struct{}{}, nil
}
