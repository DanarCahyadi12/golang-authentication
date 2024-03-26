package repository

import (
	"context"
	"errors"
	"golang-authentication/internal/entity"
	"gorm.io/gorm"
)

type UserRepositoryInterface interface {
	Save(ctx context.Context, user *entity.User) (*entity.User, error)
	DeleteById(ctx context.Context, id int) error
	FindOneByEmail(ctx context.Context, email string) (*entity.User, error)
}

type UserRepository struct {
	Database *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		Database: db,
	}
}

func (r *UserRepository) Save(ctx context.Context, user *entity.User) (*entity.User, error) {
	err := r.Database.Model(&entity.User{}).WithContext(ctx).Create(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil

}

func (r *UserRepository) DeleteById(ctx context.Context, id int) error {
	err := r.Database.WithContext(ctx).Delete(&entity.User{}, id).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) FindOneByEmail(ctx context.Context, email string) (*entity.User, error) {
	var user *entity.User
	err := r.Database.Model(&entity.User{}).WithContext(ctx).First(&user, "email = ?", email).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, err
		}

	}

	return user, nil
}
