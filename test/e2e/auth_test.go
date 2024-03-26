package e2e

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/require"
	"golang-authentication/internal/entity"
	"golang-authentication/internal/models"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAuth(t *testing.T) {
	t.Run("Should sign in and refresh token added to cookie", func(t *testing.T) {
		user := entity.User{
			Name:     "danar",
			Email:    "danar@gmail.com",
			Password: "$2a$10$rzGrygHegWythHS9wnC8u.jdM7MAgqFoUsPuTIMnIugZSWa5hsfUS",
		}
		userExists, err := UserRepository.FindOneByEmail(context.Background(), user.Email)
		require.Nil(t, err)
		if userExists != nil {
			err := UserRepository.DeleteById(context.Background(), userExists.Id)
			require.Nil(t, err)
		}
		savedUser, err := UserRepository.Save(context.Background(), &user)
		require.Nil(t, err)
		require.NotNil(t, savedUser)

		body := models.SignInRequest{
			Email:    user.Email,
			Password: "12345678",
		}

		bodyJson, err := json.Marshal(body)
		require.Nil(t, err)

		req := httptest.NewRequest(http.MethodPost, "/auth", strings.NewReader(string(bodyJson)))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")

		response, err := App.Fiber.Test(req)
		require.Nil(t, err)

		var expectedResponse models.Response[*models.UserResponse]
		bodyByte, err := io.ReadAll(response.Body)
		require.Nil(t, err)

		err = json.Unmarshal(bodyByte, &expectedResponse)
		require.Nil(t, err)

		cookieHeader := response.Header.Get("Set-Cookie")
		fmt.Println("Cookie : ", strings.Split(strings.Split(cookieHeader, ";")[0], "=")[1])
		require.NotNil(t, cookieHeader)
		require.Equal(t, "Sign in successfully", expectedResponse.Message)
		require.NotNil(t, expectedResponse.Data)
		require.Equal(t, http.StatusOK, response.StatusCode)

		err = UserRepository.DeleteById(context.Background(), user.Id)
		require.Nil(t, err)

	})
}
