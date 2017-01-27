package Messages

import (
	"github.com/gocql/gocql"
)

type Message struct {
	ID           gocql.UUID `json:"id"`
	UserID       gocql.UUID `json:"user_id"`
	UserFullName string `json:"user_full_name"`
	Message      string `json:"lastname"`
}

type GetMessageResponse struct {
	Message Message `json:"message"`
}

type MessagesResponse struct {
	Messages []Message `json:"messages"`
}

type NewMessageResponse struct {
	ID gocql.UUID `json:"id"`
}

type ErrorResponse struct {
	Errors []string `json:"errors"`
}

