package store

import (
	"github.com/satori/go.uuid"
	"github.com/KingCrunch/mocksmtp/mail"
	"time"
	"errors"
)

type memory struct {
	mailBucket []mail.Mail
}

func NewMemoryStore() (*memory) {
	return &memory{
		mailBucket: make([]mail.Mail,0,32),
	}
}

func (s *memory) Get(id uuid.UUID) (mail.Mail, error) {
	for _, el := range s.mailBucket {
		if (el.Id == id) {
			return el, nil
		}
	}
	return mail.Mail{}, errors.New("Not found")
}

func (s *memory) List() ([]mail.Mail, error) {
	return s.mailBucket, nil
}

func (s *memory) Push(mail mail.Mail) (error) {
	s.mailBucket = append(s.mailBucket, mail)

	return nil
}

func (s *memory) PurgeBefore(t time.Time) error {
	newBucket := make([]mail.Mail, 0, len(s.mailBucket))
	for _, mail := range s.mailBucket {
		if (!mail.ReceivedAt.Before(t)) {
			newBucket = append(newBucket, mail)
		}
	}
	s.mailBucket = newBucket

	return nil
}

func (s *memory) Purge() (error) {
	s.mailBucket = make([]mail.Mail, 0, len(s.mailBucket))

	return nil
}
