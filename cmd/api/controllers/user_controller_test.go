package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	errApi "chat-system/cmd/api/errors"
	"chat-system/cmd/api/models"
	mocks "chat-system/cmd/api/services/mocks"

	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRegisterUser(t *testing.T) {
	mockUserService := new(mocks.UserService)
	controller := NewUserController(mockUserService)

	// Setup Gin router and recorder
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/users/register", controller.RegisterUser)

	user := models.User{
		Name:     "John",
		Surname:  "Doe",
		Username: "johndoe",
		Password: "password",
		Email:    "johndoe@example.com",
	}

	mockUserService.On("RegisterUser", mock.Anything, user).Return(nil)

	// Marshal user to JSON
	userJSON, _ := json.Marshal(user)
	req, _ := http.NewRequest(http.MethodPost, "/users/register", bytes.NewBuffer(userJSON))
	req.Header.Set("Content-Type", "application/json")

	// Record response
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)
	mockUserService.AssertExpectations(t)
}

func TestRegisterUser_BadRequest(t *testing.T) {
	mockUserService := new(mocks.UserService)
	controller := NewUserController(mockUserService)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/users/register", controller.RegisterUser)

	invalidJSON := `{"username": "johndoe", "password": "password"`
	req, _ := http.NewRequest(http.MethodPost, "/users/register", bytes.NewBuffer([]byte(invalidJSON)))
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
	mockUserService.AssertNotCalled(t, "RegisterUser")
}

func TestRegisterUserError(t *testing.T) {
	mockUserService := new(mocks.UserService)
	controller := NewUserController(mockUserService)

	// Setup Gin router and recorder
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/users/register", controller.RegisterUser)

	user := models.User{
		Name:     "John",
		Surname:  "Doe",
		Username: "johndoe",
		Password: "password",
		Email:    "johndoe@example.com",
	}

	mockUserService.On("RegisterUser", mock.Anything, user).Return(errApi.NewErrAPIInternalServer(fmt.Errorf("internal server error")))

	// Marshal user to JSON
	userJSON, _ := json.Marshal(user)
	req, _ := http.NewRequest(http.MethodPost, "/users/register", bytes.NewBuffer(userJSON))
	req.Header.Set("Content-Type", "application/json")

	// Record response
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
	mockUserService.AssertExpectations(t)
}

func TestVerifyUserPassword(t *testing.T) {
	mockUserService := new(mocks.UserService)
	controller := NewUserController(mockUserService)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/users/login", controller.VerifyUserPassword)

	loginRequest := models.LoginRequest{
		Username: "johndoe",
		Password: "password",
	}

	mockUserService.On("VerifyUserPassword", mock.Anything, loginRequest.Username, loginRequest.Password).
		Return("1234", true, nil)

	loginJSON, _ := json.Marshal(loginRequest)
	req, _ := http.NewRequest(http.MethodPost, "/users/login", bytes.NewBuffer(loginJSON))
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)
	mockUserService.AssertExpectations(t)
}

func TestVerifyUserPassword_BadRequest(t *testing.T) {
	mockUserService := new(mocks.UserService)
	controller := NewUserController(mockUserService)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/users/login", controller.VerifyUserPassword)

	invalidJSON := `{"username": "johndoe", "password": "password"`
	req, _ := http.NewRequest(http.MethodPost, "/users/login", bytes.NewBuffer([]byte(invalidJSON)))
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
	mockUserService.AssertNotCalled(t, "RegisterUser")
}

func TestVerifyUserPassword_Unauthorized(t *testing.T) {
	mockUserService := new(mocks.UserService)
	controller := NewUserController(mockUserService)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/users/login", controller.VerifyUserPassword)

	loginRequest := models.LoginRequest{
		Username: "johndoe",
		Password: "wrongpassword",
	}

	mockUserService.On("VerifyUserPassword", mock.Anything, loginRequest.Username, loginRequest.Password).
		Return("", false, nil)
	loginJSON, _ := json.Marshal(loginRequest)
	req, _ := http.NewRequest(http.MethodPost, "/users/login", bytes.NewBuffer(loginJSON))
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusUnauthorized, recorder.Code)
	mockUserService.AssertExpectations(t)
}

func TestListUsers(t *testing.T) {
	mockUserService := new(mocks.UserService)
	controller := NewUserController(mockUserService)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/users", controller.ListUsers)

	users := []models.User{
		{
			Id:        gocql.TimeUUID(),
			Name:      "John",
			Surname:   "Doe",
			Username:  "johndoe",
			Email:     "johndoe@example.com",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	mockUserService.On("ListUsers", mock.Anything).Return(users, nil)

	req, _ := http.NewRequest(http.MethodGet, "/users", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)
	mockUserService.AssertExpectations(t)
}

func TestListUsers_NotFound(t *testing.T) {
	mockUserService := new(mocks.UserService)
	controller := NewUserController(mockUserService)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/users", controller.ListUsers)

	mockUserService.On("ListUsers", mock.Anything).
		Return(nil, errApi.NewErrAPINotFound(errors.New("no users found")))

	req, _ := http.NewRequest(http.MethodGet, "/users", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusNotFound, recorder.Code)
	mockUserService.AssertExpectations(t)
}
