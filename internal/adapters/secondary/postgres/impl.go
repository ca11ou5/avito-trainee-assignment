package postgres

import (
	"context"
	"fmt"
	"github.com/ca11ou5/avito-trainee-assignment/internal/payload"
)

const (
	isEmployeeExistsQuery  = `SELECT COUNT(*) FROM employee WHERE username = $1;`
	insertEmployeeQuery    = `INSERT INTO employee (username, hashed_password) VALUES ($1, $2);`
	getHashedPasswordQuery = `SELECT hashed_password FROM employee WHERE username = $1;`
)

func (a *Adapter) IsEmployeeExists(ctx context.Context, username string) (bool, error) {
	var count int

	err := a.db.GetContext(ctx, &count, isEmployeeExistsQuery, username)
	if err != nil {
		return false, fmt.Errorf("db get query: %s", err)
	}

	return count > 0, nil
}

func (a *Adapter) InsertEmployee(ctx context.Context, req payload.AuthRequest) error {
	_, err := a.db.ExecContext(ctx, insertEmployeeQuery, req.Username, req.Password)
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
