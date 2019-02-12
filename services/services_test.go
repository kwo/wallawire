package services_test

import (
	"context"
	"time"

	"wallawire/model"
)

type DatabaseMock struct{}

func (z *DatabaseMock) Run(fn func(tx model.Transaction) error) error {
	tx := new(TransactionMock)
	return fn(tx)
}

type TransactionMock struct{}

func (z *TransactionMock) Query(query string, params map[string]interface{}) (model.Rows, error) {
	return nil, nil
}

func (z *TransactionMock) Exec(query string, params map[string]interface{}) (model.Result, error) {
	return nil, nil
}

type UserRepositoryMock struct {
	User           *model.User
	Roles          []model.UserRole
	Available      bool
	AvailableError error
	GetError       error
	RolesError     error
	SetError       error
}

func (z *UserRepositoryMock) IsUsernameAvailable(ctx context.Context, tx model.ReadOnlyTransaction, username string) (bool, error) {
	return z.Available, z.AvailableError
}

func (z *UserRepositoryMock) GetUser(ctx context.Context, tx model.ReadOnlyTransaction, userID string) (*model.User, error) {
	return z.User, z.GetError
}

func (z *UserRepositoryMock) SetUser(ctx context.Context, tx model.WriteOnlyTransaction, user model.User) error {
	return z.SetError
}

func (z *UserRepositoryMock) GetActiveUserByUsername(ctx context.Context, tx model.ReadOnlyTransaction, username string) (*model.User, error) {
	return z.User, z.GetError
}

func (z *UserRepositoryMock) GetUserRoles(ctx context.Context, tx model.ReadOnlyTransaction, userID string, t *time.Time) ([]model.UserRole, error) {
	return z.Roles, z.RolesError
}

type IdGeneratorMock struct {
	ID string
}

func (z *IdGeneratorMock) NewID() string {
	return z.ID
}
