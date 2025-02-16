package postgres

import (
	"github.com/ca11ou5/avito-trainee-assignment/internal/models"
	"github.com/jmoiron/sqlx"

	"context"
	"errors"
	"fmt"
)

const (
	isEmployeeExistsQuery  = `SELECT 1 FROM employee WHERE username = $1;`
	insertEmployeeQuery    = `INSERT INTO employee (username, hashed_password) VALUES ($1, $2);`
	getHashedPasswordQuery = `SELECT hashed_password FROM employee WHERE username = $1;` //nolint:gosec

	getEmployeeReceiverTransactions = `SELECT sender_username, amount FROM transaction WHERE receiver_username = $1;`
	getEmployeeSentTransactions     = `SELECT receiver_username, amount FROM transaction WHERE sender_username = $1;`
	getEmployeeMerch                = `SELECT merch_name, count FROM employee_merch WHERE employee_username = $1;`
	getEmployeeBalance              = `SELECT balance FROM employee WHERE username = $1;`

	updateEmployeeBalancesQuery = `
	WITH deducted AS (
		UPDATE employee 
		SET balance = balance - $1 
		WHERE username = $2 AND balance >= $1
		RETURNING username
	), added AS (
		UPDATE employee 
		SET balance = balance + $1 
		WHERE username = $3
		RETURNING username
	)
	INSERT INTO transaction (sender_username, receiver_username, amount)
	SELECT $2, $3, $1 
	WHERE EXISTS (SELECT 1 FROM deducted) AND EXISTS (SELECT 1 FROM added);
	`

	isMerchExistsQuery = `SELECT 1 FROM merch WHERE name = $1;`

	insertEmployeeMerchQuery = `
	WITH deducted AS (
	    UPDATE employee 
	    SET balance = balance - (SELECT cost FROM merch WHERE name = $2)
	    WHERE username = $1 
	      AND balance >= (SELECT cost FROM merch WHERE name = $2)
	    RETURNING username
	), inserted AS (
	    INSERT INTO employee_merch (employee_username, merch_name, count)
	    VALUES ($1, $2, 1)
	    ON CONFLICT (employee_username, merch_name)
	    DO UPDATE SET count = employee_merch.count + 1
	    RETURNING employee_username, merch_name, count
	)
	SELECT * FROM inserted WHERE EXISTS (SELECT 1 FROM deducted);
	`
)

var (
	ErrNotEnoughBalance  = errors.New("not enough balance")
	ErrMerchNotExists    = errors.New("chosen merch does not exist")
	ErrEmployeeNotExists = errors.New("employee does not exist")
)

func (a *Adapter) IsEmployeeExists(ctx context.Context, username string) error {
	var exists bool

	err := a.db.GetContext(ctx, &exists, isEmployeeExistsQuery, username)
	if err != nil {
		return ErrEmployeeNotExists
	}

	return nil
}

func (a *Adapter) InsertEmployee(ctx context.Context, creds models.Credentials) error {
	_, err := a.db.ExecContext(ctx, insertEmployeeQuery, creds.Username, creds.Password)
	if err != nil {
		return fmt.Errorf("db exec query: %s", err)
	}

	return nil
}

func (a *Adapter) GetHashedPassword(ctx context.Context, username string) (string, error) {
	var hashedPassword string

	err := a.db.GetContext(ctx, &hashedPassword, getHashedPasswordQuery, username)
	if err != nil {
		return "", fmt.Errorf("db get query: %s", err)
	}

	return hashedPassword, nil
}

func (a *Adapter) GetEmployeeInfo(ctx context.Context, username string) (models.EmployeeInfo, error) {
	info := models.EmployeeInfo{}

	tx, err := a.db.BeginTxx(ctx, nil)
	if err != nil {
		return info, fmt.Errorf("init transaction: %s", err)
	}

	err = func() error {
		if info.CoinHistory, err = a.getCoinHistory(ctx, tx, username); err != nil {
			return fmt.Errorf("get coin history: %s", err)
		}

		if info.Inventory, err = a.getEmployeeMerch(ctx, tx, username); err != nil {
			return fmt.Errorf("get merch: %s", err)
		}

		if info.Coins, err = a.getEmployeeBalance(ctx, tx, username); err != nil {
			return fmt.Errorf("get balance: %s", err)
		}

		return nil
	}()
	if err != nil {
		err = tx.Rollback()
		if err != nil {
			return models.EmployeeInfo{}, fmt.Errorf("rollback transaction: %s", err)
		}
		return models.EmployeeInfo{}, err
	}

	return info, tx.Commit()
}

func (a *Adapter) getCoinHistory(ctx context.Context, tx *sqlx.Tx, username string) (models.CoinHistory, error) {
	var history models.CoinHistory

	recTransactions := []models.ReceivedTransaction{}
	err := tx.SelectContext(ctx, &recTransactions, getEmployeeReceiverTransactions, username)
	if err != nil {
		return history, fmt.Errorf("db select query: %s", err)
	}

	sentTransactions := []models.SentTransaction{}
	err = tx.SelectContext(ctx, &sentTransactions, getEmployeeSentTransactions, username)
	if err != nil {
		return history, fmt.Errorf("db select query: %s", err)
	}

	history = models.CoinHistory{
		Sent:     sentTransactions,
		Received: recTransactions,
	}

	return history, nil
}

func (a *Adapter) getEmployeeMerch(ctx context.Context, tx *sqlx.Tx, username string) ([]models.Item, error) {
	items := []models.Item{}

	err := tx.SelectContext(ctx, &items, getEmployeeMerch, username)
	if err != nil {
		return nil, fmt.Errorf("db select query: %s", err)
	}

	return items, nil
}

func (a *Adapter) getEmployeeBalance(ctx context.Context, tx *sqlx.Tx, username string) (int, error) {
	var balance int

	err := tx.GetContext(ctx, &balance, getEmployeeBalance, username)
	if err != nil {
		return -1, fmt.Errorf("db get query: %s", err)
	}

	return balance, nil
}

func (a *Adapter) SendCoin(ctx context.Context, username string, trans models.SentTransaction) error {
	tx, err := a.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("init transaction: %s", err)
	}

	err = func() error {
		balance, err := a.getEmployeeBalance(ctx, tx, username)
		if err != nil {
			return fmt.Errorf("get balance: %s", err)
		}

		if balance < trans.Amount {
			return ErrNotEnoughBalance
		}

		res, err := tx.ExecContext(ctx, updateEmployeeBalancesQuery, trans.Amount, username, trans.ToUser)
		if err != nil {
			return fmt.Errorf("exec send coin transaction: %s", err)
		}

		affected, err := res.RowsAffected()
		if err != nil || affected == 0 {
			return ErrNotEnoughBalance
		}

		return nil
	}()
	if err != nil {
		err = tx.Rollback()
		if err != nil {
			return fmt.Errorf("rollback transaction: %s", err)
		}
		return err
	}

	return tx.Commit()
}

func (a *Adapter) IsMerchExists(ctx context.Context, itemName string) error {
	var exists bool

	err := a.db.GetContext(ctx, &exists, isMerchExistsQuery, itemName)
	if err != nil {
		return ErrMerchNotExists
	}

	return nil
}

func (a *Adapter) InsertEmployeeMerch(ctx context.Context, username string, merch string) error {
	tx, err := a.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("init transaction: %s", err)
	}

	res, err := tx.Exec(insertEmployeeMerchQuery, username, merch)
	if err != nil {
		err = tx.Rollback()
		if err != nil {
			return fmt.Errorf("rollback transaction: %s", err)
		}
		return fmt.Errorf("exec insert merch: %s", err)
	}

	affected, err := res.RowsAffected()
	if affected == 0 || err != nil {
		tx.Rollback() //nolint:errcheck
		return ErrNotEnoughBalance
	}

	return tx.Commit()
}
