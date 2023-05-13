package event_sourcing

import (
	"github.com/google/uuid"
)

type IEventRepository interface {
	Save(event *Event) error
	Read(uuid.UUID, string) (*Event, error)
	Delete(uuid.UUID, string) error
}
