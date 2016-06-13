package store

import (
	"github.com/satori/go.uuid"
	"github.com/KingCrunch/visualsmtp/mail"
	"time"
)

type Store interface {
	Get(id uuid.UUID) (mail.Mail, error)
	List() (map[uuid.UUID]mail.Mail, error)
	Push(mail mail.Mail) (error)
	PurgeBefore(t time.Time) (error)
	Purge() (error)
}
