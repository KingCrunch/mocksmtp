package store

import (
	"github.com/satori/go.uuid"
	"github.com/KingCrunch/visualsmtp/mail"
)

type InMemory struct {
	mailBucket map[uuid.UUID]mail.Mail
}

func NewInMemoryStore () (InMemory){
	return InMemory{
		mailBucket: make(map[uuid.UUID] mail.Mail),
	}
}

func (s InMemory) Get(id uuid.UUID) (mail.Mail, error) {
	return s.mailBucket[id], nil
}

func (s InMemory) List() (map[uuid.UUID]mail.Mail, error) {
	return s.mailBucket, nil
}

func (s InMemory) Push(mail mail.Mail) (uuid.UUID, error) {
	id := uuid.NewV4()
	s.mailBucket[id] = mail

	return id, nil
}
