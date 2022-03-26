package service

import (
	"github.com/b-open/jobbuzz/pkg/model"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ErrUserAlreadyExist = errors.New("User already exist")
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (s *Service) CreateUser(email string, password string) error {
	// TODO: validate email and check password
	// create user
	passwordHash, err := HashPassword(password)
	if err != nil {
		return errors.Wrapf(err, "Error in HashPassword")
	}

	user := model.User{
		Email:        email,
		PasswordHash: passwordHash,
	}

	err = s.DB.Transaction(func(tx *gorm.DB) error {
		var existingUser *model.User
		result := tx.Find(&existingUser, model.User{Email: email})
		if result.Error != nil {
			return errors.Wrapf(result.Error, "Error getting user")
		}
		if result.RowsAffected > 0 {
			return ErrUserAlreadyExist
		}

		result = tx.Create(&user)
		if result.Error != nil {
			return errors.Wrapf(result.Error, "Error inserting user")
		}

		return nil
	})
	if err != nil {
		return errors.Wrapf(err, "Error in Transaction")
	}

	return nil
	// TODO: create jwt and return
}
