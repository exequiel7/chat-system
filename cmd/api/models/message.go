package models

import (
	"time"

	"github.com/gocql/gocql"
)

type Message struct {
	ID             gocql.UUID `json:"-" cql:"id uuid"`
	ConversationID gocql.UUID `json:"-" cql:"conversation_id uuid"`
	SenderID       gocql.UUID `json:"sender_id" cql:"sender_id uuid"`
	ReceiverID     gocql.UUID `json:"receiver_id" cql:"receiver_id uuid"`
	Content        string     `json:"content" cql:"content text"`
	Timestamp      time.Time  `json:"timestamp" cql:"timestamp timestamp"`
}
