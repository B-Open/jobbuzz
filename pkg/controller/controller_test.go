package controller

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/b-open/jobbuzz/pkg/graph/graphmodel"
	"github.com/b-open/jobbuzz/pkg/model"
	"github.com/b-open/jobbuzz/pkg/pagination"
	"github.com/gin-gonic/gin"
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

func TestGetJobs(t *testing.T) {
	t.Run("Get empty jobs", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		service := MockService{}
		service.On("GetJobs").Return([]*model.Job{}, &pagination.PaginationResult{}, nil)

		controller := Controller{Service: &service}
		controller.GetJobs(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("return 1 job", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		service := MockService{}
		service.On("GetJobs").Return([]*model.Job{
			{
				BaseModel: model.BaseModel{
					ID: 1,
				},
				Title: "test job",
			},
		}, &pagination.PaginationResult{}, nil)

		controller := Controller{Service: &service}
		controller.GetJobs(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("error should panic", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		service := MockService{}
		service.On("GetJobs").Return([]*model.Job{}, &pagination.PaginationResult{}, errors.New("service error"))

		controller := Controller{Service: &service}

		assert.PanicsWithError(t, "service error", func() { controller.GetJobs(c) })
	})
}
