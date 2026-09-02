package rest

import (
	"context"
	"encoding/json"
	"math"
	"net/http"

	"github.com/mytheresa/go-hiring-challenge/domain/list"
	"github.com/mytheresa/go-hiring-challenge/domain/query"
	"github.com/mytheresa/go-hiring-challenge/lib/operation"
	"github.com/mytheresa/go-hiring-challenge/shared"
)

type ListResponse[T any] struct {
	Metadata ListResponseMetadata `json:"metadata"`
	Data     []T                  `json:"data"`
}

type ListResponseMetadata struct {
	OperationId string `json:"operationId"`
	TotalCount  uint   `json:"totalCount"`
	Page        int    `json:"page"`
	PageSize    int    `json:"pageSize"`
	PageCount   int    `json:"pageCount"`
}

func EncodeListResponse[DO, RO any](
	log shared.Logger,
	ctx context.Context,
	itemEncoder func(DO) RO,
	out list.ListResponse[DO],
	query *query.Query,
	w http.ResponseWriter,
	r *http.Request,
) {
	pageCount := math.Ceil(float64(out.Total()) / float64((query.Pagination().PageSize())))

	restOut := ListResponse[RO]{
		Metadata: ListResponseMetadata{
			OperationId: operation.IdFromContext(ctx),
			TotalCount:  out.Total(),
			Page:        query.Pagination().Page(),
			PageSize:    query.Pagination().PageSize(),
			PageCount:   int(pageCount),
		},
		Data: []RO{},
	}

	if len(out.Items()) > 0 {
		for i := range out.Items() {
			restOut.Data = append(restOut.Data, itemEncoder(out.Items()[i]))
		}
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
