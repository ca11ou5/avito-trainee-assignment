package postgres

import (
	"context"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ca11ou5/avito-trainee-assignment/internal/models"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

func TestInsertEmployee(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")
	adapter := &Adapter{
		db: sqlxDB,
	}

	creds := models.Credentials{
		Username: "testuser",
		Password: "hashedPassword",
	}

	mock.ExpectExec(regexp.QuoteMeta(insertEmployeeQuery)).
		WithArgs(creds.Username, creds.Password).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = adapter.InsertEmployee(context.Background(), creds)
	assert.NoError(t, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestIsEmployeeExists_EmployeeNotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")
	adapter := &Adapter{
		db: sqlxDB,
	}

	username := "nonexistent_user"
	mock.ExpectQuery(regexp.QuoteMeta(isEmployeeExistsQuery)).
		WithArgs(username).
		WillReturnError(fmt.Errorf("no rows in result set"))

	err = adapter.IsEmployeeExists(context.Background(), username)
	assert.Equal(t, ErrEmployeeNotExists, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestGetHashedPassword(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")
	adapter := &Adapter{
		db: sqlxDB,
	}

	username := "testuser"
	expectedHash := "somehashedvalue"

	rows := sqlmock.NewRows([]string{"hashed_password"}).AddRow(expectedHash)
	mock.ExpectQuery(regexp.QuoteMeta(getHashedPasswordQuery)).
		WithArgs(username).
		WillReturnRows(rows)

	actualHash, err := adapter.GetHashedPassword(context.Background(), username)
	assert.NoError(t, err)
	assert.Equal(t, expectedHash, actualHash)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestSendCoin_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")
	adapter := &Adapter{
		db: sqlxDB,
	}

	username := "testuser"
	transAmount := 100
	toUser := "receiver"

	// Ожидаем начало транзакции
	mock.ExpectBegin()

	// Ожидаем вызов getEmployeeBalance внутри транзакции.
	// Здесь запрос getEmployeeBalance ожидается с нужным аргументом.
	mock.ExpectQuery(regexp.QuoteMeta(getEmployeeBalance)).
		WithArgs(username).
		WillReturnRows(sqlmock.NewRows([]string{"balance"}).AddRow(200))

	mock.ExpectExec(regexp.QuoteMeta(updateEmployeeBalancesQuery)).
		WithArgs(transAmount, username, toUser).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	err = adapter.SendCoin(context.Background(), username, models.SentTransaction{
		Amount: transAmount,
		ToUser: toUser,
	})
	assert.NoError(t, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}
