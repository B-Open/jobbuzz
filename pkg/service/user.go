package service

import (
	"fmt"
	"time"

	"github.com/b-open/jobbuzz/pkg/model"
	"github.com/golang-jwt/jwt/v4"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ErrUserAlreadyExist = errors.New("User already exist")
)

type (
	Claims struct {
		Email string `json:"email"`
		jwt.RegisteredClaims
	}
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func generateAccessToken(user model.User) (string, error) {
	claims := Claims{
		user.Email,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "https://jobbuzz.org",
			Subject:   fmt.Sprintf("%d", user.ID),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte("mysupersecretkey")) // TODO: use config
	if err != nil {
		return "", errors.Wrapf(err, "Error in SignedString")
	}

	return ss, nil
}

func (s *Service) CreateUser(email string, password string) (token string, err error) {
	// TODO: validate email and check password
	// create user
	passwordHash, err := HashPassword(password)
	if err != nil {
		return "", errors.Wrapf(err, "Error in HashPassword")
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
		if err == ErrUserAlreadyExist {
			return "", err
		}
		return "", errors.Wrapf(err, "Error in Transaction")
	}

	token, err = generateAccessToken(user)
	if err != nil {
		return "", errors.Wrapf(err, "Error in generateAccessToken")
	}

	return token, nil
}
