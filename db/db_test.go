package db

import (
	"database/sql"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestUserExists_WhenExists(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %v", err)
	}
	defer mockDB.Close()

	// replace package DB with mock
	DB = mockDB

	rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id FROM user WHERE username = ?")).WithArgs("alice").WillReturnRows(rows)

	exists, err := UserExists("alice")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !exists {
		t.Fatalf("expected user to exist")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestUserExists_WhenNotExists(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %v", err)
	}
	defer mockDB.Close()

	DB = mockDB

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id FROM user WHERE username = ?")).WithArgs("bob").WillReturnError(sql.ErrNoRows)

	exists, err := UserExists("bob")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if exists {
		t.Fatalf("expected user to NOT exist")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestCreateUser_Success(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %v", err)
	}
	defer mockDB.Close()

	DB = mockDB

	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO user (username, password) VALUES (?, ?)")).WithArgs("carol", "hashed").WillReturnResult(sqlmock.NewResult(5, 1))

	if err := CreateUser("carol", "hashed"); err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}
