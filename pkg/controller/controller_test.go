package controller

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/b-open/jobbuzz/pkg/model"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockService struct {
	mock.Mock
}

func (s *MockService) GetJobs() ([]model.Job, error) {
	args := s.Called()
	return args.Get(0).([]model.Job), args.Error(1)
}

func TestGetJobs(t *testing.T) {
	t.Run("Get empty jobs", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		service := MockService{}
		service.On("GetJobs").Return([]model.Job{}, nil)

		controller := Controller{Service: &service}
		controller.GetJobs(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("return 1 job", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		service := MockService{}
		service.On("GetJobs").Return([]model.Job{
			{
				BaseModel: model.BaseModel{
					ID: 1,
				},
				Title: "test job",
			},
		}, nil)

		controller := Controller{Service: &service}
		controller.GetJobs(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("error should panic", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		service := MockService{}
		service.On("GetJobs").Return([]model.Job{}, errors.New("service error"))

		controller := Controller{Service: &service}

		assert.PanicsWithError(t, "service error", func() { controller.GetJobs(c) })
	})
}
