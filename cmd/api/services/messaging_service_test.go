package services

import (
	"context"
	"errors"
	"testing"
	"time"

	errApi "chat-system/cmd/api/errors"
	"chat-system/cmd/api/models"
	"chat-system/cmd/api/repositories/mocks"

	"github.com/gocql/gocql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSendMessage(t *testing.T) {
	mockRepo := new(mocks.DBRepository)
	mockService := NewMessagingService(mockRepo)

	msg := models.Message{
		SenderID:   gocql.TimeUUID(),
		ReceiverID: gocql.TimeUUID(),
		Content:    "Hello",
	}

	mockRepo.On("SaveMessage", mock.Anything, mock.Anything).Return(nil)

	err := mockService.SendMessage(context.Background(), msg)
	assert.Nil(t, err)

	mockRepo.AssertExpectations(t)
}

func TestSendMessage_EmptyContent(t *testing.T) {
	mockRepo := new(mocks.DBRepository)
	mockService := NewMessagingService(mockRepo)

	msg := models.Message{
		SenderID:   gocql.TimeUUID(),
		ReceiverID: gocql.TimeUUID(),
		Content:    "",
	}

	err := mockService.SendMessage(context.Background(), msg)
	assert.NotNil(t, err)
	assert.Equal(t, errApi.NewErrAPIBadRequest(errors.New("message content cannot be empty")).GetMessage(), err.GetMessage())
}

func TestSendMessage_SaveError(t *testing.T) {
	mockRepo := new(mocks.DBRepository)
	mockService := NewMessagingService(mockRepo)

	msg := models.Message{
		SenderID:   gocql.TimeUUID(),
		ReceiverID: gocql.TimeUUID(),
		Content:    "Hello",
	}

	mockRepo.On("SaveMessage", mock.Anything, mock.Anything).Return(errors.New("some error"))

	err := mockService.SendMessage(context.Background(), msg)
	assert.NotNil(t, err)
	assert.Equal(t, errApi.NewErrAPIInternalServer(errors.New("some error")).GetMessage(), err.GetMessage())
}

func TestGetConversationHistory(t *testing.T) {
	mockRepo := new(mocks.DBRepository)
	mockService := NewMessagingService(mockRepo)

	userID1 := gocql.TimeUUID()
	userID2 := gocql.TimeUUID()
	conversationID := generateConversationID(userID1, userID2)

	expectedMessages := []models.Message{
		{
			ID:             gocql.TimeUUID(),
			ConversationID: conversationID,
			SenderID:       userID1,
			ReceiverID:     userID2,
			Content:        "Hello",
			Timestamp:      time.Now(),
		},
	}

	mockRepo.On("GetMessagesBetweenUsers", mock.Anything, conversationID).Return(expectedMessages, nil)

	messages, err := mockService.GetConversationHistory(context.Background(), userID1, userID2)
	assert.Nil(t, err)
	assert.Equal(t, expectedMessages, messages)

	mockRepo.AssertExpectations(t)
}

func TestGetConversationHistory_NotFound(t *testing.T) {
	mockRepo := new(mocks.DBRepository)
	mockService := NewMessagingService(mockRepo)

	userID1 := gocql.TimeUUID()
	userID2 := gocql.TimeUUID()
	conversationID := generateConversationID(userID1, userID2)

	mockRepo.On("GetMessagesBetweenUsers", mock.Anything, conversationID).Return(nil, nil)

	messages, err := mockService.GetConversationHistory(context.Background(), userID1, userID2)
	assert.NotNil(t, err)
	assert.Nil(t, messages)
	assert.Equal(t, errApi.NewErrAPINotFound(errors.New("no messages found between the specified users")).GetMessage(), err.GetMessage())

	mockRepo.AssertExpectations(t)
}

func TestGetConversationHistory_Error(t *testing.T) {
	mockRepo := new(mocks.DBRepository)
	mockService := NewMessagingService(mockRepo)

	userID1 := gocql.TimeUUID()
	userID2 := gocql.TimeUUID()
	conversationID := generateConversationID(userID1, userID2)

	mockRepo.On("GetMessagesBetweenUsers", mock.Anything, conversationID).Return(nil, errors.New("some error"))

	messages, err := mockService.GetConversationHistory(context.Background(), userID1, userID2)
	assert.NotNil(t, err)
	assert.Nil(t, messages)
	assert.Equal(t, errApi.NewErrAPIInternalServer(errors.New("some error")).GetMessage(), err.GetMessage())

	mockRepo.AssertExpectations(t)
}
