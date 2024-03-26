package e2e

import (
	"context"
	"encoding/json"
	"github.com/stretchr/testify/require"
	"golang-authentication/internal/entity"
	"golang-authentication/internal/models"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSignUp(t *testing.T) {
	t.Run("Empty field name", func(t *testing.T) {
		body := models.SignUpRequest{
			Name:     "",
			Email:    "danar@gmail.com",
			Password: "12345",
		}
		bodyJson, err := json.Marshal(body)
		require.Nil(t, err)
		req := httptest.NewRequest(http.MethodPost, "/signup", strings.NewReader(string(bodyJson)))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")

		response, err := App.Fiber.Test(req)
		require.Nil(t, err)

		var expectedResponse models.ErrorResponse
		bodyByte, err := io.ReadAll(response.Body)
		require.Nil(t, err)

		err = json.Unmarshal(bodyByte, &expectedResponse)
		require.Nil(t, err)

		require.Equal(t, "Name must be required", expectedResponse.Message)
		require.Equal(t, "Bad Request", expectedResponse.Status)
		require.Equal(t, http.StatusBadRequest, response.StatusCode)
	})

	t.Run("Empty field email", func(t *testing.T) {
		body := models.SignUpRequest{
			Name:     "Danar",
			Email:    "",
			Password: "12345",
		}
		bodyJson, err := json.Marshal(body)
		require.Nil(t, err)
		req := httptest.NewRequest(http.MethodPost, "/signup", strings.NewReader(string(bodyJson)))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")

		response, err := App.Fiber.Test(req)
		require.Nil(t, err)

		var expectedResponse models.ErrorResponse
		bodyByte, err := io.ReadAll(response.Body)
		require.Nil(t, err)

		err = json.Unmarshal(bodyByte, &expectedResponse)
		require.Nil(t, err)

		require.Equal(t, "Email must be required", expectedResponse.Message)
		require.Equal(t, "Bad Request", expectedResponse.Status)
		require.Equal(t, http.StatusBadRequest, response.StatusCode)
	})

	t.Run("Empty field password", func(t *testing.T) {
		body := models.SignUpRequest{
			Name:     "Danar",
			Email:    "danar@gmail.com",
			Password: "",
		}
		bodyJson, err := json.Marshal(body)
		require.Nil(t, err)
		req := httptest.NewRequest(http.MethodPost, "/signup", strings.NewReader(string(bodyJson)))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")

		response, err := App.Fiber.Test(req)
		require.Nil(t, err)

		var expectedResponse models.ErrorResponse
		bodyByte, err := io.ReadAll(response.Body)
		require.Nil(t, err)

		err = json.Unmarshal(bodyByte, &expectedResponse)
		require.Nil(t, err)

		require.Equal(t, "Password must be required", expectedResponse.Message)
		require.Equal(t, "Bad Request", expectedResponse.Status)
		require.Equal(t, http.StatusBadRequest, response.StatusCode)
	})
	t.Run("Password less 8 than 8 character", func(t *testing.T) {
		body := models.SignUpRequest{
			Name:     "Danar",
			Email:    "danar@gmail.com",
			Password: "12345",
		}
		bodyJson, err := json.Marshal(body)
		require.Nil(t, err)
		req := httptest.NewRequest(http.MethodPost, "/signup", strings.NewReader(string(bodyJson)))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")

		response, err := App.Fiber.Test(req)
		require.Nil(t, err)

		var expectedResponse models.ErrorResponse
		bodyByte, err := io.ReadAll(response.Body)
		require.Nil(t, err)

		err = json.Unmarshal(bodyByte, &expectedResponse)
		require.Nil(t, err)

		require.Equal(t, "Password must be min 8", expectedResponse.Message)
		require.Equal(t, "Bad Request", expectedResponse.Status)
		require.Equal(t, http.StatusBadRequest, response.StatusCode)
	})

	t.Run("Not email", func(t *testing.T) {
		body := models.SignUpRequest{
			Name:     "Danar",
			Email:    "gmail",
			Password: "12345768",
		}
		bodyJson, err := json.Marshal(body)
		require.Nil(t, err)
		req := httptest.NewRequest(http.MethodPost, "/signup", strings.NewReader(string(bodyJson)))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")

		response, err := App.Fiber.Test(req)
		require.Nil(t, err)

		var expectedResponse models.ErrorResponse
		bodyByte, err := io.ReadAll(response.Body)
		require.Nil(t, err)

		err = json.Unmarshal(bodyByte, &expectedResponse)
		require.Nil(t, err)

		require.Equal(t, "Email must be email", expectedResponse.Message)
		require.Equal(t, "Bad Request", expectedResponse.Status)
		require.Equal(t, http.StatusBadRequest, response.StatusCode)
	})

	t.Run("Unique email", func(t *testing.T) {
		user := entity.User{
			Name:     "danar",
			Email:    "danar@gmail.com",
			Password: "12345678",
		}
		savedUser, err := UserRepository.Save(context.Background(), &user)
		require.Nil(t, err)
		require.NotNil(t, savedUser)

		body := models.SignUpRequest{
			Name:     "Danar",
			Email:    "danar@gmail.com",
			Password: "12345768",
		}
		bodyJson, err := json.Marshal(body)
		require.Nil(t, err)

		req := httptest.NewRequest(http.MethodPost, "/signup", strings.NewReader(string(bodyJson)))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")

		response, err := App.Fiber.Test(req)
		require.Nil(t, err)

		var expectedResponse models.ErrorResponse
		bodyByte, err := io.ReadAll(response.Body)
		require.Nil(t, err)

		err = json.Unmarshal(bodyByte, &expectedResponse)
		require.Nil(t, err)

		require.Equal(t, "Email already exists", expectedResponse.Message)
		require.Equal(t, "Bad Request", expectedResponse.Status)
		require.Equal(t, http.StatusBadRequest, response.StatusCode)

		err = UserRepository.DeleteById(context.Background(), savedUser.Id)
		require.Nil(t, err)
	})

	t.Run("Correct signup", func(t *testing.T) {
		body := models.SignUpRequest{
			Name:     "Danar",
			Email:    "danar@gmail.com",
			Password: "12345768",
		}
		bodyJson, err := json.Marshal(body)
		require.Nil(t, err)
		req := httptest.NewRequest(http.MethodPost, "/signup", strings.NewReader(string(bodyJson)))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")

		response, err := App.Fiber.Test(req)
		require.Nil(t, err)

		var expectedResponse models.Response[*models.UserResponse]
		bodyByte, err := io.ReadAll(response.Body)
		require.Nil(t, err)

		err = json.Unmarshal(bodyByte, &expectedResponse)
		require.Nil(t, err)

		require.Equal(t, "Signup successfully", expectedResponse.Message)
		require.NotNil(t, expectedResponse.Data)
		require.Equal(t, http.StatusCreated, response.StatusCode)

		err = UserRepository.DeleteById(context.Background(), expectedResponse.Data.Id)
		require.Nil(t, err)
	})

}
