package mocks

import (
	"context"
	"github.com/stretchr/testify/mock"
	"golang-authentication/internal/entity"
)

type UserRepositoryMock struct {
	Mock mock.Mock
}

func NewUserRepositoryMock() *UserRepositoryMock {
	return &UserRepositoryMock{
		Mock: mock.Mock{},
	}
}

func (r *UserRepositoryMock) Save(ctx context.Context, user *entity.User) (*entity.User, error) {
	args := r.Mock.Called(user)
	return args.Get(0).(*entity.User), nil

}

func (r *UserRepositoryMock) DeleteById(ctx context.Context, id int) error {
	args := r.Mock.Called(id)
	return args.Error(0)

}

func (r *UserRepositoryMock) FindOneByEmail(ctx context.Context, email string) (*entity.User, error) {
	args := r.Mock.Called(email)
	if args.Get(0) == nil {
		return nil, nil
	}
	return args.Get(0).(*entity.User), nil
}
