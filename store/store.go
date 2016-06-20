package store

import (
	"github.com/satori/go.uuid"
	"time"
	"github.com/veqryn/go-email/email"
)

type Item struct {
	Id uuid.UUID
	Sender string
	Recipients []string
	Message *email.Message
	FirstSeen time.Time
}

func NewItem(message *email.Message, sender string, recipients []string) (*Item) {
	return &Item {
		Id: uuid.NewV4(),
		Sender: sender,
		Recipients: recipients,
		Message: message,
		FirstSeen: time.Now(),
	}
}

type Store interface {
	Get(id uuid.UUID) (*Item, error)
	List() ([]*Item, error)
	Push(item *Item) (error)
	PurgeBefore(t time.Time) (error)
	Purge() (error)
}
