package service

import (
	"github.com/b-open/jobbuzz/pkg/model"
)

// TODO: rename file to user.go
func (s *Service) CreateUser(email string, password string) error {
	// TODO: validate email and check password
	// create user
	// TODO: hash and salt password
	user := model.User{
		Email: email,
	}

	result := s.DB.FirstOrCreate(&user, model.User{Email: email})
	if result.Error != nil {
		return result.Error
	}

	return nil
	// TODO: create jwt and return
}
