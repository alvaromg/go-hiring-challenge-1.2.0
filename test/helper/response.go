package helper

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

type ListMetadataDTO struct {
	OperationId string `json:"operationId"`
	TotalCount  uint   `json:"totalCount"`
	Page        int    `json:"page"`
	PageSize    int    `json:"pageSize"`
	PageCount   int    `json:"pageCount"`
}

type ListResponseDTO[T any] struct {
	Metadata ListMetadataDTO `json:"metadata"`
	Data     T               `json:"data"`
}

type ResponseDTO[T any] struct {
	Data T `json:"data"`
}

type ErrorResponseDTO struct {
	Error struct {
		Message string `json:"message"`
	} `json:"error"`
}

func DecodeBody(t *testing.T, rec *httptest.ResponseRecorder, out any) {
	t.Helper()
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), out))
}
