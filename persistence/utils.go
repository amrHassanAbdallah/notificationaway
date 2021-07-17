package persistence

import (
	"go.mongodb.org/mongo-driver/bson"
)

// QueryNotifierFilters defines model for QueryNotifierFilters.
type QueryMessageFilters struct {
	Ids []string `json:"ids,omitempty"`
}

// QueryNotifiersOptions defines model for QueryNotifiersOptions.
type QueryMessagesOptions struct {
	Filters       QueryMessageFilters `json:"filters"`
	Limit         int64               `json:"limit" validate:"max=1000"`
	Page          int64               `json:"page" validate:"min=0"`
	SortBy        string              `json:"sort_by" validate:"oneof=created_at updated_at name"`
	SortDirection string              `json:"sort_direction" validate:"oneof=asc desc"`
	Id            string              `json:"id"`
}

//NotifiersPaginationOptions
type QueryNotifiersPagination struct {
	Limit         int64  `json:"limit"`
	Page          int64  `json:"page"`
	SortBy        string `json:"sort_by"`
	SortDirection int    `json:"sort_direction"`
}

func constructQueryNotifiers(payload QueryMessagesOptions) (bson.D, QueryNotifiersPagination) {
	filters := bson.D{
		{Key: "deleted_at", Value: bson.M{"$exists": false}},
	}

	ops := QueryNotifiersPagination{
		Limit:  payload.Limit,
		Page:   payload.Page,
		SortBy: payload.SortBy,
	}

	if payload.Id != "" {
		filters = append(filters, bson.E{Key: "_id", Value: payload.Id})
	}

	if payload.Filters.Ids != nil {
		filters = append(filters, bson.E{Key: "_id", Value: bson.M{"$in": payload.Filters.Ids}})
	}

	if payload.SortDirection == "asc" {
		ops.SortDirection = 1
	} else {
		ops.SortDirection = -1
	}

	return filters, ops
}
