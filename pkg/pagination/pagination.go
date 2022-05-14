package pagination

import (
	"math"

	"github.com/b-open/jobbuzz/pkg/graph/graphmodel"
)

type (
	PaginationResult struct {
		To          int `json:"to"`
		From        int `json:"from"`
		PerPage     int `json:"per_page"`
		CurrentPage int `json:"current_page"`
		TotalPage   int `json:"total_page"`
		Total       int `json:"total"`
	}
)

func GetOffset(input graphmodel.PaginationInput) (offset int) {
	return (input.Page - 1) * input.Size
}

func GetPaginationResult(input graphmodel.PaginationInput, dataLength int, totalCount int) PaginationResult {
	result := PaginationResult{
		To:          GetOffset(input) + dataLength,
		From:        GetOffset(input) + 1,
		PerPage:     input.Size,
		CurrentPage: input.Page,
		TotalPage:   int(math.Ceil(float64(totalCount) / float64(input.Size))),
		Total:       dataLength,
	}

	return result
}
