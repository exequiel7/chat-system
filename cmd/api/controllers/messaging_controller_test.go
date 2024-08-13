package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
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

func TestSendMessage(t *testing.T) {
	mockMessagingService := new(mocks.MessagingService)
	controller := NewMessagingController(mockMessagingService)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/messages/send", controller.SendMessage)

	message := models.Message{
		SenderID:   gocql.TimeUUID(),
		ReceiverID: gocql.TimeUUID(),
		Content:    "Hello, World!",
		Timestamp:  time.Now(),
	}

	mockMessagingService.On("SendMessage", mock.Anything, mock.AnythingOfType("models.Message")).Return(nil)

	messageJSON, _ := json.Marshal(message)
	req, _ := http.NewRequest(http.MethodPost, "/messages/send", bytes.NewBuffer(messageJSON))
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)
	mockMessagingService.AssertExpectations(t)
}

func TestSendMessage_BadRequest(t *testing.T) {
	mockMessagingService := new(mocks.MessagingService)
	controller := NewMessagingController(mockMessagingService)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/messages/send", controller.SendMessage)

	invalidJSON := `{"sender_id": "invalid-uuid", "content": "Hello, World!"}`
	req, _ := http.NewRequest(http.MethodPost, "/messages/send", bytes.NewBuffer([]byte(invalidJSON)))
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
	mockMessagingService.AssertNotCalled(t, "SendMessage")
}

func TestGetConversationHistory(t *testing.T) {
	mockMessagingService := new(mocks.MessagingService)
	controller := NewMessagingController(mockMessagingService)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/messages/history/:senderID/:receiverID", controller.GetConversationHistory)

	senderUUID := gocql.TimeUUID()
	receiverUUID := gocql.TimeUUID()

	messages := []models.Message{
		{
			ID:             gocql.TimeUUID(),
			ConversationID: gocql.TimeUUID(),
			SenderID:       senderUUID,
			ReceiverID:     receiverUUID,
			Content:        "Hello!",
			Timestamp:      time.Now(),
		},
	}

	mockMessagingService.On("GetConversationHistory", mock.Anything, senderUUID, receiverUUID).Return(messages, nil)

	req, _ := http.NewRequest(http.MethodGet, "/messages/history/"+senderUUID.String()+"/"+receiverUUID.String(), nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)
	mockMessagingService.AssertExpectations(t)
}

func TestGetConversationHistory_InvalidUserID(t *testing.T) {
	mockMessagingService := new(mocks.MessagingService)
	controller := NewMessagingController(mockMessagingService)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/messages/history/:senderID/:receiverID", controller.GetConversationHistory)

	req, _ := http.NewRequest(http.MethodGet, "/messages/history/invalid-sender-id/invalid-receiver-id", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
	mockMessagingService.AssertNotCalled(t, "GetConversationHistory")
}

func TestGetConversationHistory_NotFound(t *testing.T) {
	mockMessagingService := new(mocks.MessagingService)
	controller := NewMessagingController(mockMessagingService)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/messages/history/:senderID/:receiverID", controller.GetConversationHistory)

	senderUUID := gocql.TimeUUID()
	receiverUUID := gocql.TimeUUID()

	mockMessagingService.On("GetConversationHistory", mock.Anything, senderUUID, receiverUUID).
		Return(nil, errApi.NewErrAPINotFound(errors.New("no messages found between the specified users")))

	req, _ := http.NewRequest(http.MethodGet, "/messages/history/"+senderUUID.String()+"/"+receiverUUID.String(), nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusNotFound, recorder.Code)
	mockMessagingService.AssertExpectations(t)
}
