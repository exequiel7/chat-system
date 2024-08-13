package repositories

import (
	errApi "chat-system/cmd/api/errors"
	"chat-system/cmd/api/models"
	"context"
	"fmt"

	"github.com/gocql/gocql"
)

type DBRepository interface {
	SaveUser(ctx context.Context, user models.User) error
	GetUserPassword(ctx context.Context, username string) (userID string, storedHash string, errs errApi.ErrAPI)
	UserExists(ctx context.Context, username string) (bool, error)
	SaveMessage(ctx context.Context, message models.Message) error
	GetMessagesBetweenUsers(ctx context.Context, conversationID gocql.UUID) ([]models.Message, error)
	GetAllUsers(ctx context.Context) ([]models.User, error)
}

type dbRepositoryImpl struct {
	session *gocql.Session
}

func NewDBRepository(session *gocql.Session) DBRepository {
	return &dbRepositoryImpl{session: session}
}

func (r *dbRepositoryImpl) SaveUser(ctx context.Context, user models.User) error {
	query := `INSERT INTO users (id, name, surname, username, password, email, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	return r.session.Query(query, user.Id, user.Name, user.Surname, user.Username, user.Password, user.Email, user.CreatedAt, user.UpdatedAt).WithContext(ctx).Exec()
}

func (r *dbRepositoryImpl) GetUserPassword(ctx context.Context, username string) (userID string, storedHash string, errs errApi.ErrAPI) {
	query := `SELECT id, password FROM users WHERE username = ?`
	err := r.session.Query(query, username).WithContext(ctx).Scan(&userID, &storedHash)

	if err == gocql.ErrNotFound {
		return "", "", errApi.NewErrAPINotFound(fmt.Errorf("user not found"))
	}

	if err != nil {
		return "", "", errApi.NewErrAPIInternalServer(err)
	}

	return userID, storedHash, nil
}

func (r *dbRepositoryImpl) UserExists(ctx context.Context, username string) (bool, error) {
	var id gocql.UUID
	query := `SELECT id FROM users WHERE username = ? LIMIT 1`
	err := r.session.Query(query, username).WithContext(ctx).Scan(&id)
	if err == gocql.ErrNotFound {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *dbRepositoryImpl) SaveMessage(ctx context.Context, message models.Message) error {
	query := `INSERT INTO messages (id, conversation_id, sender_id, receiver_id, content, timestamp) VALUES (?, ?, ?, ?, ?, ?)`
	return r.session.Query(query, message.ID, message.ConversationID, message.SenderID, message.ReceiverID, message.Content, message.Timestamp).WithContext(ctx).Exec()
}

func (r *dbRepositoryImpl) GetMessagesBetweenUsers(ctx context.Context, conversationID gocql.UUID) ([]models.Message, error) {
	var messages []models.Message
	query := `SELECT id, conversation_id, sender_id, receiver_id, content, timestamp FROM messages 
			  WHERE conversation_id = ?
			  ORDER BY timestamp ASC`

	iter := r.session.Query(query, conversationID).WithContext(ctx).Iter()

	var msg models.Message
	for iter.Scan(&msg.ID, &msg.ConversationID, &msg.SenderID, &msg.ReceiverID, &msg.Content, &msg.Timestamp) {
		messages = append(messages, msg)
	}

	if err := iter.Close(); err != nil {
		return nil, err
	}

	return messages, nil
}

func (r *dbRepositoryImpl) GetAllUsers(ctx context.Context) ([]models.User, error) {
	var users []models.User

	query := `SELECT id, name, surname, username, email, created_at, updated_at FROM users`
	iter := r.session.Query(query).WithContext(ctx).Iter()

	var user models.User
	for iter.Scan(&user.Id, &user.Name, &user.Surname, &user.Username, &user.Email, &user.CreatedAt, &user.UpdatedAt) {
		users = append(users, user)
	}

	if err := iter.Close(); err != nil {
		return nil, err
	}

	return users, nil
}
