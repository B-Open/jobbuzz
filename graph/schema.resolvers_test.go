package graph

import (
	"testing"

	"github.com/b-open/jobbuzz/pkg/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockService struct {
	mock.Mock
}

func (s *MockService) GetJobs() ([]*model.Job, error) {
	args := s.Called()
	return args.Get(0).([]*model.Job), args.Error(1)
}

func TestJobs(t *testing.T) {
	mockService := MockService{}
	mockService.On("GetJobs").Return([]*model.Job{
		{
			BaseModel: model.BaseModel{
				ID: 1,
			},
			Title: "test job",
		},
	}, nil)

	r := Resolver{Service: &mockService}

	result, err := r.Query().Jobs(nil)
	if err != nil {
		t.Fatal(err)
	}

	assert.NotEmpty(t, result, "Jobs should not be empty")
}
