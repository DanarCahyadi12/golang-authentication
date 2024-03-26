package usecase

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"golang-authentication/internal/entity"
	"golang-authentication/internal/helpers"
	"golang-authentication/internal/models"
	"golang-authentication/internal/repository"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type SignUpUseCase struct {
	UserRepository repository.UserRepositoryInterface
	Validator      *validator.Validate
}

func NewSignupUseCase(userRepository repository.UserRepositoryInterface, validator *validator.Validate) *SignUpUseCase {
	return &SignUpUseCase{UserRepository: userRepository, Validator: validator}
}

func (u *SignUpUseCase) CreateUser(ctx context.Context, userRequest *models.SignUpRequest) (*models.UserResponse, error) {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err := u.validateRequest(ctxWithTimeout, userRequest)
	if err != nil {

		return nil, err
	}

	hashedPassword, err := u.HashPassword(userRequest.Password)
	if err != nil {
		return nil, err
	}

	user := &entity.User{
		Name:     userRequest.Name,
		Email:    userRequest.Email,
		Password: hashedPassword,
	}
	result, err := u.UserRepository.Save(ctxWithTimeout, user)

	if errors.Is(err, context.DeadlineExceeded) {
		return nil, &models.ErrorResponse{Code: 408, Message: "Request timeout. Please try again", Status: "Request Timeout"}
	}
	return &models.UserResponse{
		Id:        result.Id,
		Name:      result.Name,
		Email:     result.Email,
		CreatedAt: result.CreatedAt,
		UpdatedAt: result.UpdatedAt,
	}, nil
}

func (u *SignUpUseCase) HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("Error while hashing password :", err)
		return "", &models.ErrorResponse{Code: 500, Message: "Something wrong", Status: "Internal Server Error"}
	}

	return string(hashedPassword), nil

}

func (u *SignUpUseCase) validateRequest(ctx context.Context, userRequest *models.SignUpRequest) error {

	err := u.Validator.Struct(userRequest)
	if err != nil {
		if e := err.(validator.ValidationErrors); e != nil {
			message := helpers.GetFirstValidationErrorsAndConvert(err)
			fmt.Println("Error while validating: ", e)
			return &models.ErrorResponse{
				Message: message,
				Status:  "Bad Request",
				Code:    400,
			}
		}
		fmt.Println("Server error while validating: ", err)
		return &models.ErrorResponse{Message: "Something wrong", Code: 500, Status: "Internal Server Error"}

	}

	//if email is already taken
	userExist, err := u.UserRepository.FindOneByEmail(ctx, userRequest.Email)
	if err != nil {
		fmt.Println("Something error while getting user by email: ", err)
		if errors.Is(err, context.DeadlineExceeded) {
			return &models.ErrorResponse{Code: 408, Status: "Request Timeout", Message: "Request timeout. Please try again"}
		}
		return &models.ErrorResponse{Message: "Something wrong!", Code: 500, Status: "Internal Server Error"}

	}

	if userExist != nil {
		return &models.ErrorResponse{Message: "Email already exists", Code: 400, Status: "Bad Request"}
	}
	return nil
}
