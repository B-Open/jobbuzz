package graph

import (
	"errors"
	"testing"

	"github.com/b-open/jobbuzz/pkg/graph/graphmodel"
	"github.com/b-open/jobbuzz/pkg/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockService struct {
	mock.Mock
}

func (s *MockService) GetJobs(pagination graphmodel.PaginationInput) ([]*model.Job, int64, error) {
	args := s.Called()

	jobs := args.Get(0)
	if jobs == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}

	return args.Get(0).([]*model.Job), args.Get(1).(int64), args.Error(2)
}

func (s *MockService) GetCompanies(pagination graphmodel.PaginationInput) ([]*model.Company, int64, error) {
	args := s.Called()

	jobs := args.Get(0)
	if jobs == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}

	return args.Get(0).([]*model.Company), args.Get(1).(int64), args.Error(2)
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
		}, int64(1), nil)

		r := Resolver{Service: &mockService}

		result, err := r.Query().Jobs(nil, nil, graphmodel.PaginationInput{})
		if err != nil {
			t.Fatal(err)
		}

		assert.NotEmpty(t, result.Data, "Jobs should not be empty")
		assert.Len(t, result.Data, 1)
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
		mockService.On("GetJobs").Return(mockJobs, int64(20), nil)

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
		mockService.On("GetJobs").Return([]*model.Job{}, int64(0), nil)

		r := Resolver{Service: &mockService}

		result, err := r.Query().Jobs(nil, nil, graphmodel.PaginationInput{})
		if err != nil {
			t.Fatal(err)
		}

		assert.Empty(t, result.Data, "Jobs should be empty")
	})

	t.Run("test error", func(t *testing.T) {
		mockService := MockService{}
		mockService.On("GetJobs").Return(nil, int64(0), errors.New("error"))

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
		}, int64(1), nil)

		r := Resolver{Service: &mockService}

		got, err := r.Query().Companies(nil, nil, graphmodel.PaginationInput{})

		assert.Nil(t, err)
		assert.NotEmpty(t, got)
	})
}
