package services

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	errApi "chat-system/cmd/api/errors"
	"chat-system/cmd/api/models"
	dbMocks "chat-system/cmd/api/repositories/mocks"

	"github.com/gocql/gocql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRegisterUser(t *testing.T) {
	mockRepo := new(dbMocks.DBRepository)
	userService := NewUserService(mockRepo)

	user := models.User{
		Id:        gocql.TimeUUID(),
		Name:      "John",
		Surname:   "Doe",
		Username:  "johndoe",
		Password:  "password",
		Email:     "johndoe@example.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockRepo.On("UserExists", mock.Anything, user.Username).Return(false, nil)
	mockRepo.On("SaveUser", mock.Anything, mock.Anything).Return(nil)

	err := userService.RegisterUser(context.Background(), user)
	assert.Nil(t, err)

	mockRepo.AssertExpectations(t)
}

func TestRegisterUser_UserAlreadyExists(t *testing.T) {
	mockRepo := new(dbMocks.DBRepository)
	userService := NewUserService(mockRepo)

	user := models.User{
		Username: "johndoe",
	}

	mockRepo.On("UserExists", mock.Anything, user.Username).Return(true, nil)

	err := userService.RegisterUser(context.Background(), user)
	assert.NotNil(t, err)
	assert.Contains(t, err.GetMessage(), "username already exists")

	mockRepo.AssertExpectations(t)
}

func TestRegisterUser_SaveError(t *testing.T) {
	mockRepo := new(dbMocks.DBRepository)
	userService := NewUserService(mockRepo)

	user := models.User{
		Id:        gocql.TimeUUID(),
		Name:      "John",
		Surname:   "Doe",
		Username:  "johndoe",
		Password:  "password",
		Email:     "johndoe@example.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockRepo.On("UserExists", mock.Anything, user.Username).Return(false, nil)
	mockRepo.On("SaveUser", mock.Anything, mock.Anything).Return(errors.New("some error"))

	err := userService.RegisterUser(context.Background(), user)
	assert.NotNil(t, err)

	mockRepo.AssertExpectations(t)
}

func TestVerifyUserPassword(t *testing.T) {
	mockRepo := new(dbMocks.DBRepository)
	UserService := NewUserService(mockRepo)

	storedHash := "hashedpassword"

	mockRepo.On("GetUserPassword", mock.Anything, "johndoe").Return("1234", storedHash, nil)

	_, hash, err := UserService.VerifyUserPassword(context.Background(), "johndoe", "password")
	assert.Nil(t, err)
	assert.NotNil(t, hash)

	mockRepo.AssertExpectations(t)
}

func TestVerifyUserPassword_InvalidCredentials(t *testing.T) {
	mockRepo := new(dbMocks.DBRepository)
	userService := NewUserService(mockRepo)

	storedHash := "hashedpassword"

	mockRepo.On("GetUserPassword", mock.Anything, "johndoe").Return("1234", storedHash, nil)

	_, success, err := userService.VerifyUserPassword(context.Background(), "johndoe", "wrongpassword")
	assert.Nil(t, err)
	assert.False(t, success)

	mockRepo.AssertExpectations(t)
}

func TestVerifyUserPassword_UserNotFound(t *testing.T) {
	mockRepo := new(dbMocks.DBRepository)
	userService := NewUserService(mockRepo)

	mockRepo.On("GetUserPassword", mock.Anything, "johndoe").Return("", "", errApi.NewErrAPIBadRequest(fmt.Errorf("User not found")))

	_, success, err := userService.VerifyUserPassword(context.Background(), "johndoe", "password")
	assert.NotNil(t, err)
	assert.False(t, success)

	mockRepo.AssertExpectations(t)
}

func TestListUsers(t *testing.T) {
	mockRepo := new(dbMocks.DBRepository)
	userService := NewUserService(mockRepo)

	expectedUsers := []models.User{
		{
			Id:        gocql.TimeUUID(),
			Name:      "John",
			Surname:   "Doe",
			Username:  "johndoe",
			Password:  "password",
			Email:     "johndoe@example.com",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			Id:        gocql.TimeUUID(),
			Name:      "John",
			Surname:   "Doe",
			Username:  "johndoe2",
			Password:  "password",
			Email:     "johndoe2@example.com",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	mockRepo.On("GetAllUsers", mock.Anything).Return(expectedUsers, nil)

	users, err := userService.ListUsers(context.Background())
	assert.Nil(t, err)
	assert.Equal(t, expectedUsers, users)

	mockRepo.AssertExpectations(t)
}

func TestListUsers_Error(t *testing.T) {
	mockRepo := new(dbMocks.DBRepository)
	userService := NewUserService(mockRepo)

	mockRepo.On("GetAllUsers", mock.Anything).Return(nil, errors.New("some error"))

	users, err := userService.ListUsers(context.Background())
	assert.NotNil(t, err)
	assert.Nil(t, users)

	mockRepo.AssertExpectations(t)
}
