package usecase

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
	"golang-authentication/internal/entity"
	"golang-authentication/internal/helpers"
	"golang-authentication/internal/models"
	"golang-authentication/internal/repository"
	"golang.org/x/crypto/bcrypt"
	"sync"
	"time"
)

type AuthUseCase struct {
	UserRepository repository.UserRepositoryInterface
	Validator      *validator.Validate
	Viper          *viper.Viper
}

func NewAuthUseCase(userRepository repository.UserRepositoryInterface, validator *validator.Validate, viper *viper.Viper) *AuthUseCase {
	return &AuthUseCase{
		UserRepository: userRepository,
		Validator:      validator,
		Viper:          viper,
	}
}
func (u *AuthUseCase) GenerateAccessToken(userID int) (string, error) {
	key := u.Viper.GetString("key.token.access")

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": "restful-api",
		"sub": userID,
		"exp": time.Now().Add(1 * time.Hour).Unix(),
	})

	token, err := jwtToken.SignedString([]byte(key))
	if err != nil {
		fmt.Println("Error while generate access token : ", err)
		return "", &models.ErrorResponse{
			Code:    500,
			Message: "Error while generate access token",
			Status:  "Internal Server Error",
		}
	}

	return token, nil

}

func (u *AuthUseCase) GenerateRefreshToken(userID int) (string, error) {
	key := u.Viper.GetString("key.token.refresh")

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": "restful-api",
		"sub": userID,
		"exp": time.Now().Add(3 * (24 * time.Hour)).Unix(),
	})

	token, err := jwtToken.SignedString([]byte(key))
	if err != nil {
		fmt.Println("Error while generate access token : ", err)
		return "", &models.ErrorResponse{
			Code:    500,
			Message: "Error while generate access token",
			Status:  "Internal Server Error",
		}
	}

	return token, nil
}
func (u *AuthUseCase) ValidateRequest(request *models.SignInRequest) error {
	err := u.Validator.Struct(request)
	if err != nil {
		if e := err.(validator.ValidationErrors); e != nil {
			message := helpers.GetFirstValidationErrorsAndConvert(err)
			return &models.ErrorResponse{Code: 400, Message: message, Status: "Bad Request"}
		}

		return &models.ErrorResponse{Code: 500, Message: "Something wrong", Status: "Internal Server Error"}
	}

	return nil

}

func (u *AuthUseCase) GetAndValidateUser(ctx context.Context, credential *models.SignInRequest) (*entity.User, error) {
	user, err := u.UserRepository.FindOneByEmail(ctx, credential.Email)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, &models.ErrorResponse{Code: 408, Status: "Request Timeout", Message: "Request timeout. Please try again"}
		}
		return nil, &models.ErrorResponse{Code: 500, Status: "Internal Server Error", Message: "Something wrong!"}
	}

	if user == nil {
		return nil, &models.ErrorResponse{Code: 400, Status: "Bad Request", Message: "Email or password invalid"}
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credential.Password))
	if err != nil {
		return nil, &models.ErrorResponse{Code: 400, Status: "Bad Request", Message: "Email or password invalid"}
	}

	return user, nil

}
func (u *AuthUseCase) SignIn(ctx context.Context, credential *models.SignInRequest) (*models.SignInResponse, error) {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	err := u.ValidateRequest(credential)
	if err != nil {
		fmt.Println("Error while validating request: ", err)
		return nil, err
	}

	user, err := u.GetAndValidateUser(ctxWithTimeout, credential)
	if err != nil {
		return nil, err
	}
	userID := user.Id

	var accessToken string
	var refreshToken string
	var wg sync.WaitGroup
	var errorChannel = make(chan error, 2)

	//generate tokens using goroutine
	wg.Add(2)
	go func() {
		defer wg.Done()
		accessToken, err = u.GenerateAccessToken(userID)
		if err != nil {
			errorChannel <- err
			return
		}

		errorChannel <- nil

	}()
	go func() {
		defer wg.Done()
		refreshToken, err = u.GenerateRefreshToken(userID)
		if err != nil {
			errorChannel <- err
			return
		}

		errorChannel <- nil

	}()

	//wait until both goroutine finished the process
	wg.Wait()

	//if errorChannel have a error, then return the error
	for i := 0; i < 2; i++ {
		if err := <-errorChannel; err != nil {
			return nil, err
		}
	}
	close(errorChannel)

	return &models.SignInResponse{AccessToken: accessToken, RefreshToken: refreshToken}, nil

}

func (u *AuthUseCase) VerifyRefreshToken(refreshToken string, key string) (float64, error) {
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})

	if err != nil {
		fmt.Println("Error while parsing token, ", err)
		return -1, &models.ErrorResponse{
			Code:    401,
			Message: "Invalid token",
			Status:  "Unauthorized",
		}
	}

	if !token.Valid {
		return -1, &models.ErrorResponse{
			Code:    401,
			Message: "Invalid token",
			Status:  "Unauthorized",
		}
	}
	var sub float64
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		sub = claims["sub"].(float64)
	}

	return sub, nil
}

func (u *AuthUseCase) GetToken(refreshToken string) (*models.GetTokenResponse, error) {
	if refreshToken == "" {
		return nil, &models.ErrorResponse{
			Code:    401,
			Message: "Please sign in first",
			Status:  "Unauthorized",
		}
	}
	refreshTokenKey := u.Viper.GetString("key.token.refresh")
	sub, err := u.VerifyRefreshToken(refreshToken, refreshTokenKey)
	if err != nil {
		return nil, err
	}

	accessToken, err := u.GenerateAccessToken(int(sub))
	if err != nil {
		return nil, err
	}

	return &models.GetTokenResponse{AccessToken: accessToken}, nil

}
