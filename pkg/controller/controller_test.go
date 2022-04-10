package controller

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/b-open/jobbuzz/pkg/graph/graphmodel"
	"github.com/b-open/jobbuzz/pkg/model"
	"github.com/gin-gonic/gin"
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

func TestGetJobs(t *testing.T) {
	t.Run("Get empty jobs", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		service := MockService{}
		service.On("GetJobs").Return([]*model.Job{}, int64(0), nil)

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
		}, int64(1), nil)

		controller := Controller{Service: &service}
		controller.GetJobs(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("error should panic", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		service := MockService{}
		service.On("GetJobs").Return([]*model.Job{}, int64(0), errors.New("service error"))

		controller := Controller{Service: &service}

		assert.PanicsWithError(t, "service error", func() { controller.GetJobs(c) })
	})
}
