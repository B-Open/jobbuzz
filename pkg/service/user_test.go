package service

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestHashPassword(t *testing.T) {
	password := "verystrongpassword123"
	passwordHash, err := HashPassword(password)

	assert.Nil(t, err)

	want := true
	got := CheckPasswordHash(password, passwordHash)
	assert.Equal(t, got, want)
}

func TestCheckPasswordHash(t *testing.T) {
	t.Run("incorrect password", func(t *testing.T) {
		password := "verystrongpassword123"
		passwordHash, err := HashPassword(password)
		passwordAttempt := "hunter2"

		assert.Nil(t, err)

		want := false
		got := CheckPasswordHash(passwordAttempt, passwordHash)
		assert.Equal(t, got, want)
	})
}

func TestCreateAccount(t *testing.T) {
	t.Run("new account", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatal(err)
		}
		defer db.Close()

		// gorm mysql version check
		mock.ExpectQuery("SELECT VERSION\\(\\)").WillReturnRows(sqlmock.NewRows([]string{"version()"}).AddRow("8.0.28"))
		mock.ExpectBegin()
		// user exist check
		mock.ExpectQuery("SELECT \\* FROM `users`").WithArgs("test@example.com").WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "email", "password"}))
		// insert new user
		mock.ExpectExec("INSERT INTO `users`").WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), "test@example.com", sqlmock.AnyArg()).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		gormDB, err := gorm.Open(mysql.New(mysql.Config{Conn: db}), &gorm.Config{})
		if err != nil {
			t.Fatal(err)
		}

		s := Service{DB: gormDB}
		err = s.CreateUser("test@example.com", "mypassword")

		assert.Nil(t, err)
		assert.Nil(t, mock.ExpectationsWereMet(), "sqlmock expectations were not met")
	})

	t.Run("account already exist should return error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatal(err)
		}
		defer db.Close()

		// gorm mysql version check
		mock.ExpectQuery("SELECT VERSION\\(\\)").WillReturnRows(sqlmock.NewRows([]string{"version()"}).AddRow("8.0.28"))
		mock.ExpectBegin()
		// user exist check
		mock.ExpectQuery("SELECT \\* FROM `users`").WithArgs("test@example.com").WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "email", "password"}).AddRow(1, nil, nil, nil, "test@example.com", ""))

		gormDB, err := gorm.Open(mysql.New(mysql.Config{Conn: db}), &gorm.Config{})
		if err != nil {
			t.Fatal(err)
		}

		s := Service{DB: gormDB}
		err = s.CreateUser("test@example.com", "mypassword")

		assert.ErrorIs(t, err, ErrUserAlreadyExist)
		assert.Nil(t, mock.ExpectationsWereMet(), "sqlmock expectations were not met")
	})
}
