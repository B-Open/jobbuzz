package service

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestCreateAccount(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	// gorm mysql version check
	mock.ExpectQuery("SELECT VERSION\\(\\)").WillReturnRows(sqlmock.NewRows([]string{"version()"}).AddRow("8.0.28"))
	// user exist check
	mock.ExpectQuery("SELECT \\* FROM `users`").WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "email", "password"}))
	// insert new user
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `users`").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	gormDB, err := gorm.Open(mysql.New(mysql.Config{Conn: db}), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}

	s := Service{DB: gormDB}
	err = s.CreateUser("test@example.com", "mypassword")

	assert.Nil(t, err)
	assert.Nil(t, mock.ExpectationsWereMet(), "sqlmock expectations were not met")
}
