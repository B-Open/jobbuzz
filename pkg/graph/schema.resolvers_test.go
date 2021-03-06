package graph

import (
	"errors"
	"testing"

	"github.com/b-open/jobbuzz/pkg/graph/graphmodel"
	"github.com/b-open/jobbuzz/pkg/model"
	"github.com/b-open/jobbuzz/pkg/pagination"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockService struct {
	mock.Mock
}

func (s *MockService) GetJobs(paginationInput graphmodel.PaginationInput) ([]*model.Job, *pagination.PaginationResult, error) {
	args := s.Called()

	jobs := args.Get(0)
	if jobs == nil {
		return nil, nil, args.Error(2)
	}

	return args.Get(0).([]*model.Job), args.Get(1).(*pagination.PaginationResult), args.Error(2)
}

func (s *MockService) GetCompanies(paginationInput graphmodel.PaginationInput) ([]*model.Company, *pagination.PaginationResult, error) {
	args := s.Called()

	jobs := args.Get(0)
	if jobs == nil {
		return nil, args.Get(1).(*pagination.PaginationResult), args.Error(2)
	}

	return args.Get(0).([]*model.Company), args.Get(1).(*pagination.PaginationResult), args.Error(2)
}

func (s *MockService) CreateUser(email string, password string) (string, error) {
	args := s.Called()
	return args.Get(0).(string), args.Error(1)
}

func TestRegisterAccount(t *testing.T) {
	mockService := MockService{}
	mockService.On("CreateUser").Return("accesstoken", nil)

	r := Resolver{Service: &mockService}

	_, err := r.Mutation().RegisterAccount(nil, graphmodel.NewUserInput{Email: "test@example.com", Password: "password"})

	assert.Nil(t, err)
}

func TestJobs(t *testing.T) {
	t.Run("test return 1 job", func(t *testing.T) {
		mockService := MockService{}
		mockService.On("GetJobs").Return([]*model.Job{
			{
				BaseModel: model.BaseModel{
					ID: 1,
				},
				Title: "test job",
			},
		}, &pagination.PaginationResult{
			To:          1,
			From:        1,
			PerPage:     10,
			CurrentPage: 1,
			TotalPage:   1,
			Total:       1,
		}, nil)

		r := Resolver{Service: &mockService}

		want := &graphmodel.JobOutput{
			To:          1,
			From:        1,
			PerPage:     10,
			CurrentPage: 1,
			TotalPage:   1,
			Total:       1,
			Data: []*graphmodel.Job{
				{
					ID:    1,
					Title: "test job",
				},
			},
		}

		got, err := r.Query().Jobs(nil, nil, graphmodel.PaginationInput{Page: 1, Size: 10})
		if err != nil {
			t.Fatal(err)
		}

		assert.NotEmpty(t, got.Data, "Jobs should not be empty")
		assert.Len(t, got.Data, 1)
		assert.Equal(t, want, got)
	})

	t.Run("test return 20 jobs", func(t *testing.T) {
		mockService := MockService{}
		var mockJobs []*model.Job
		for i := 0; i < 20; i++ {
			mockJobs = append(mockJobs, &model.Job{
				BaseModel: model.BaseModel{
					ID: uint64(i),
				},
				Title: "test job",
			})
		}
		mockService.On("GetJobs").Return(mockJobs, &pagination.PaginationResult{}, nil)

		r := Resolver{Service: &mockService}

		result, err := r.Query().Jobs(nil, nil, graphmodel.PaginationInput{})
		if err != nil {
			t.Fatal(err)
		}

		assert.NotEmpty(t, result.Data, "Jobs should not be empty")
		assert.Len(t, result.Data, 20, "Jobs length is not correct")
	})

	t.Run("test return no jobs", func(t *testing.T) {
		mockService := MockService{}
		mockService.On("GetJobs").Return([]*model.Job{}, &pagination.PaginationResult{}, nil)

		r := Resolver{Service: &mockService}

		result, err := r.Query().Jobs(nil, nil, graphmodel.PaginationInput{})
		if err != nil {
			t.Fatal(err)
		}

		assert.Empty(t, result.Data, "Jobs should be empty")
	})

	t.Run("test error", func(t *testing.T) {
		mockService := MockService{}
		mockService.On("GetJobs").Return(nil, &pagination.PaginationResult{}, errors.New("error"))

		r := Resolver{Service: &mockService}

		_, err := r.Query().Jobs(nil, nil, graphmodel.PaginationInput{})
		assert.NotNil(t, err, "Error was expected but not found.")
	})
}

func TestCompanies(t *testing.T) {
	t.Run("get 1 company", func(t *testing.T) {
		mockService := MockService{}
		mockService.On("GetCompanies").Return([]*model.Company{
			{BaseModel: model.BaseModel{ID: 1}},
		}, &pagination.PaginationResult{}, nil)

		r := Resolver{Service: &mockService}

		got, err := r.Query().Companies(nil, nil, graphmodel.PaginationInput{})

		assert.Nil(t, err)
		assert.NotEmpty(t, got)
	})
}
