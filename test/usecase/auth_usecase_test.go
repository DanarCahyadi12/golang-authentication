package usecase

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	"golang-authentication/internal/config"
	"golang-authentication/internal/entity"
	"golang-authentication/internal/models"
	"golang-authentication/internal/usecase"
	"golang-authentication/test/mocks"
	"sync"
	"testing"
)

func TestAuthUseCase(t *testing.T) {
	viper := config.NewViper("./../../")
	validator := config.NewValidator()
	repositoryMock := mocks.NewUserRepositoryMock()
	authUseCase := usecase.NewAuthUseCase(repositoryMock, validator, viper)

	t.Run("Validate request", func(t *testing.T) {
		t.Run("Sign in with empty email", func(t *testing.T) {
			model := &models.SignInRequest{
				Email:    "",
				Password: "12345678",
			}
			err := authUseCase.ValidateRequest(model)

			require.NotNil(t, err)
		})

		t.Run("Sign in with empty password", func(t *testing.T) {
			model := &models.SignInRequest{
				Email:    "danar@gmail.com",
				Password: "",
			}
			err := authUseCase.ValidateRequest(model)
			require.NotNil(t, err)

		})

	})

	t.Run("Validate user", func(t *testing.T) {
		t.Run("When user doesnt match", func(t *testing.T) {

			model := &models.SignInRequest{
				Email:    "notmatch@gmail.com",
				Password: "12345678",
			}
			repositoryMock.Mock.On("FindOneByEmail", model.Email).Return(nil)
			user, err := authUseCase.GetAndValidateUser(context.Background(), model)
			require.NotNil(t, err)
			require.Nil(t, user)
		})

		t.Run("When password user doesn't matched", func(t *testing.T) {
			user := &entity.User{
				Id:       1,
				Email:    "danar@gmail.com",
				Password: "$2a$10$rzGrygHegWythHS9wnC8u.jdM7MAgqFoUsPuTIMnIugZSWa5hsfUS",
			}
			model := &models.SignInRequest{
				Email:    "danar@gmail.com",
				Password: "1234567890",
			}

			repositoryMock.Mock.On("FindOneByEmail", model.Email).Return(user)
			user, err := authUseCase.GetAndValidateUser(context.Background(), model)
			require.Equal(t, err, &models.ErrorResponse{Code: 400, Message: "Email or password invalid", Status: "Bad Request"})
			require.Nil(t, user)

		})

		t.Run("User matched and password matched", func(t *testing.T) {
			user := &entity.User{
				Id:       1,
				Email:    "danar@gmail.com",
				Password: "$2a$10$rzGrygHegWythHS9wnC8u.jdM7MAgqFoUsPuTIMnIugZSWa5hsfUS",
			}
			model := &models.SignInRequest{
				Email:    "danar@gmail.com",
				Password: "12345678",
			}

			repositoryMock.Mock.On("FindOneByEmail", model.Email).Return(user)
			result, err := authUseCase.GetAndValidateUser(context.Background(), model)
			require.Nil(t, err)
			require.Equal(t, user, result)
		})

	})

	t.Run("Token", func(t *testing.T) {
		t.Run("Generate access token", func(t *testing.T) {
			const userID = 1
			accessToken, err := authUseCase.GenerateAccessToken(userID)
			require.Nil(t, err)
			require.NotNil(t, accessToken)

		})

		t.Run("Should running both generate token function using goroutine", func(t *testing.T) {
			var wg sync.WaitGroup
			wg.Add(2)
			go func() {
				defer wg.Done()
				const userID = 1
				accessToken, err := authUseCase.GenerateAccessToken(userID)
				require.Nil(t, err)
				require.NotNil(t, accessToken)
			}()

			go func() {
				defer wg.Done()
				const userID = 1
				refreshToken, err := authUseCase.GenerateRefreshToken(userID)
				require.Nil(t, err)
				require.NotNil(t, refreshToken)
			}()
			wg.Wait()

			fmt.Println("Token generated successfully")
		})

		t.Run("Verify refresh token should not return an error", func(t *testing.T) {
			refreshTokenKey := viper.GetString("key.token.refresh")
			require.NotNil(t, refreshTokenKey)
			refreshToken, err := authUseCase.GenerateRefreshToken(2)

			require.Nil(t, err)
			require.NotNil(t, refreshToken)

			sub, err := authUseCase.VerifyRefreshToken(refreshToken, refreshTokenKey)
			var expectedResult float64 = 2
			require.Equal(t, expectedResult, sub)
			require.Nil(t, err)
		})
		t.Run("Should generate new access token", func(t *testing.T) {
			refreshToken, err := authUseCase.GenerateRefreshToken(2)

			require.Nil(t, err)
			require.NotNil(t, refreshToken)

			result, err := authUseCase.GetToken(refreshToken)
			require.Nil(t, err)
			require.NotNil(t, result.AccessToken)
		})
	})

	t.Run("Should return access token and refresh token after sign in", func(t *testing.T) {
		user := &entity.User{
			Id:       1,
			Email:    "danar@gmail.com",
			Password: "$2a$10$rzGrygHegWythHS9wnC8u.jdM7MAgqFoUsPuTIMnIugZSWa5hsfUS",
		}
		model := &models.SignInRequest{
			Email:    "danar@gmail.com",
			Password: "12345678",
		}

		repositoryMock.Mock.On("FindOneByEmail", model.Email).Return(user)
		response, err := authUseCase.SignIn(context.Background(), model)
		require.Nil(t, err)
		require.NotNil(t, response)

	})
}
