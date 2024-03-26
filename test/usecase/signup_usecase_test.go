package usecase

import (
	"github.com/stretchr/testify/require"
	"golang-authentication/internal/config"
	"golang-authentication/internal/usecase"
	"golang-authentication/test/mocks"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func TestSignUpUseCase(t *testing.T) {
	validator := config.NewValidator()
	repositoryMock := mocks.NewUserRepositoryMock()
	signupUseCase := usecase.NewSignupUseCase(repositoryMock, validator)

	t.Run("Generate hash password and compare", func(t *testing.T) {
		hashedPassword, err := signupUseCase.HashPassword("password")
		require.Nil(t, err)
		require.NotNil(t, hashedPassword)
		err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte("password"))
		require.Nil(t, err)

	})

}
