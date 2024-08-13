package services

import (
	errApi "chat-system/cmd/api/errors"
	"chat-system/cmd/api/models"
	"chat-system/cmd/api/repositories"
	"context"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"time"

	"github.com/gocql/gocql"
)

type MessagingService interface {
	SendMessage(ctx context.Context, msg models.Message) errApi.ErrAPI
	GetConversationHistory(ctx context.Context, userID1, userID2 gocql.UUID) ([]models.Message, errApi.ErrAPI)
}

type messagingServiceImpl struct {
	dbRepository repositories.DBRepository
}

func NewMessagingService(dbRepository repositories.DBRepository) MessagingService {
	return &messagingServiceImpl{dbRepository: dbRepository}
}

func (m *messagingServiceImpl) SendMessage(ctx context.Context, msg models.Message) errApi.ErrAPI {
	if msg.Content == "" {
		return errApi.NewErrAPIBadRequest(errors.New("message content cannot be empty"))
	}

	msg.ConversationID = generateConversationID(msg.SenderID, msg.ReceiverID)
	msg.ID = gocql.TimeUUID()
	msg.Timestamp = time.Now()

	err := m.dbRepository.SaveMessage(ctx, msg)
	if err != nil {
		return errApi.NewErrAPIInternalServer(err)
	}

	return nil
}

func (m *messagingServiceImpl) GetConversationHistory(ctx context.Context, userID1, userID2 gocql.UUID) ([]models.Message, errApi.ErrAPI) {
	conversationID := generateConversationID(userID1, userID2)
	messages, err := m.dbRepository.GetMessagesBetweenUsers(ctx, conversationID)
	if err != nil {
		return nil, errApi.NewErrAPIInternalServer(err)
	}

	if len(messages) == 0 {
		return nil, errApi.NewErrAPINotFound(errors.New("no messages found between the specified users"))
	}

	return messages, nil
}

func generateConversationID(userID1, userID2 gocql.UUID) gocql.UUID {
	var id1, id2 string
	if userID1.String() < userID2.String() {
		id1 = userID1.String()
		id2 = userID2.String()
	} else {
		id1 = userID2.String()
		id2 = userID1.String()
	}

	hash := sha1.New()
	hash.Write([]byte(id1 + id2))
	sum := hash.Sum(nil)

	sha1Hex := hex.EncodeToString(sum)
	conversationID, _ := gocql.UUIDFromBytes([]byte(sha1Hex[:16]))

	return conversationID
}
