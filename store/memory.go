package store

import (
	"github.com/satori/go.uuid"
	"time"
	"errors"
)

type memory struct {
	mailBucket []*Item
}

func NewMemoryStore() (*memory) {
	return &memory{
		mailBucket: make([]*Item,0,32),
	}
}

func (s *memory) Get(id uuid.UUID) (*Item, error) {
	for _, el := range s.mailBucket {
		if (el.Id == id) {
			return el, nil
		}
	}
	return &Item{}, errors.New("Not found")
}

func (s *memory) List() ([]*Item, error) {
	return s.mailBucket, nil
}

func (s *memory) Push(item *Item) (error) {
	item.Message.Header.ContentDisposition()
	s.mailBucket = append(s.mailBucket, item)

	return nil
}

func (s *memory) PurgeBefore(t time.Time) error {
	newBucket := make([]*Item, 0, len(s.mailBucket))
	for _, item := range s.mailBucket {
		if (!item.FirstSeen.Before(t)) {
			newBucket = append(newBucket, item)
		}
	}
	s.mailBucket = newBucket

	return nil
}

func (s *memory) Purge() (error) {
	s.mailBucket = make([]*Item, 0, len(s.mailBucket))

	return nil
}
