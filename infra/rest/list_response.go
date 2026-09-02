package rest

import (
	"context"
	"encoding/json"
	"math"
	"net/http"

	"github.com/mytheresa/go-hiring-challenge/domain/list"
	"github.com/mytheresa/go-hiring-challenge/domain/query"
	"github.com/mytheresa/go-hiring-challenge/infra/operation"
	"github.com/mytheresa/go-hiring-challenge/shared"
)

type ListResponse[T any] struct {
	Metadata ListResponseMetadata `json:"metadata"`
	Data     T                    `json:"data"`
}

type ListResponseMetadata struct {
	OperationId string `json:"operationId"`
	TotalCount  uint   `json:"totalCount"`
	Page        int    `json:"page"`
	PageSize    int    `json:"pageSize"`
	PageCount   int    `json:"pageCount"`
}

// EncodeListResponse encodes a paginated domain list into the HTTP response.
// Each item is encoded via itemEncoder, then the resulting items are wrapped
// into a data envelope (e.g. {"products": [...]}) via dataWrapper.
func EncodeListResponse[DO, RO, WD any](
	log shared.Logger,
	ctx context.Context,
	itemEncoder func(DO) RO,
	dataWrapper func([]RO) WD,
	out list.ListResponse[DO],
	query *query.Query,
	w http.ResponseWriter,
	r *http.Request,
) {
	pageCount := math.Ceil(float64(out.Total()) / float64((query.Pagination().PageSize())))

	items := make([]RO, 0, len(out.Items()))
	for i := range out.Items() {
		items = append(items, itemEncoder(out.Items()[i]))
	}

	restOut := ListResponse[WD]{
		Metadata: ListResponseMetadata{
			OperationId: operation.IdFromContext(ctx),
			TotalCount:  out.Total(),
			Page:        query.Pagination().Page(),
			PageSize:    query.Pagination().PageSize(),
			PageCount:   int(pageCount),
		},
		Data: dataWrapper(items),
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
