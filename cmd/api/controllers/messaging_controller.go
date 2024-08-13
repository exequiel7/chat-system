package controllers

import (
	errApi "chat-system/cmd/api/errors"
	"chat-system/cmd/api/models"
	"chat-system/cmd/api/services"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
)

type MessagingController interface {
	SendMessage(c *gin.Context)
	GetConversationHistory(c *gin.Context)
}

type messagingControllerImpl struct {
	msgService services.MessagingService
}

func NewMessagingController(msgService services.MessagingService) MessagingController {
	return &messagingControllerImpl{msgService: msgService}
}

// SendMessage godoc
// @Summary Send a message
// @Description Sends a message from one user to another
// @Tags messages
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer Token"
// @Param message body models.Message true "Message Data"
// @Success 200 {object} models.APIResponse "Message sent successfully"
// @Failure 400 {object} errApi.ErrAPI "Invalid request payload"
// @Failure 500 {object} errApi.ErrAPI "Internal Server Error"
// @Router /messages/send [post]
func (m *messagingControllerImpl) SendMessage(c *gin.Context) {
	var msg models.Message
	if err := c.ShouldBindJSON(&msg); err != nil {
		apiErr := errApi.NewErrAPIBadRequest(errors.New("invalid request payload"))
		c.JSON(apiErr.GetHTTPStatusCode(), apiErr)
		return
	}

	err := m.msgService.SendMessage(c.Request.Context(), msg)
	if err != nil {
		c.JSON(err.GetHTTPStatusCode(), err)
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{Message: "Message sent successfully"})
}

// GetConversationHistory godoc
// @Summary Get conversation history
// @Description Retrieves the message history between two users
// @Tags messages
// @Accept json
// @Produce json
// @Param senderID path string true "Sender User ID"
// @Param receiverID path string true "Receiver User ID"
// @Success 200 {object} models.APIResponse{data=[]models.Message} "Conversation history retrieved successfully"
// @Param Authorization header string true "Bearer Token"
// @Failure 400 {object} errApi.ErrAPI "Invalid user ID"
// @Failure 404 {object} errApi.ErrAPI "No messages found between the specified users"
// @Failure 500 {object} errApi.ErrAPI "Internal Server Error"
// @Router /messages/history/{senderID}/{receiverID} [get]
func (m *messagingControllerImpl) GetConversationHistory(c *gin.Context) {
	senderID := c.Param("senderID")
	receiverID := c.Param("receiverID")

	senderUUID, err := gocql.ParseUUID(senderID)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{Message: "Invalid sender ID"})
		return
	}

	receiverUUID, err := gocql.ParseUUID(receiverID)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{Message: "Invalid receiver ID"})
		return
	}

	messages, errAPI := m.msgService.GetConversationHistory(c.Request.Context(), senderUUID, receiverUUID)
	if errAPI != nil {
		c.JSON(errAPI.GetHTTPStatusCode(), errAPI)
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Data: messages,
	})
}
