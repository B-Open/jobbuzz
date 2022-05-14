package pagination

import (
	"testing"

	"github.com/b-open/jobbuzz/pkg/graph/graphmodel"
	"github.com/stretchr/testify/assert"
)

func TestGetOffset(t *testing.T) {
	t.Run("first page", func(t *testing.T) {
		input := graphmodel.PaginationInput{
			Size: 10,
			Page: 1,
		}

		want := 0

		got := GetOffset(input)

		assert.Equal(t, want, got)
	})

	t.Run("some page", func(t *testing.T) {
		input := graphmodel.PaginationInput{
			Size: 10,
			Page: 2,
		}

		want := 10

		got := GetOffset(input)

		assert.Equal(t, want, got)
	})
}

func TestGetPaginationResult(t *testing.T) {
	t.Run("page 1 result", func(t *testing.T) {
		input := graphmodel.PaginationInput{
			Size: 10,
			Page: 1,
		}
		dataLength := 5
		totalCount := 5

		want := PaginationResult{
			To:          5,
			From:        1,
			PerPage:     10,
			CurrentPage: 1,
			TotalPage:   1,
			Total:       5,
		}

		got := GetPaginationResult(input, dataLength, totalCount)

		assert.Equal(t, want, got)
	})

	t.Run("page 1 with multiple pages", func(t *testing.T) {
		input := graphmodel.PaginationInput{
			Size: 10,
			Page: 1,
		}
		dataLength := 10
		totalCount := 33

		want := PaginationResult{
			To:          10,
			From:        1,
			PerPage:     10,
			CurrentPage: 1,
			TotalPage:   4,
			Total:       10,
		}

		got := GetPaginationResult(input, dataLength, totalCount)

		assert.Equal(t, want, got)
	})
}
