package event_sourcing

import (
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

type EventService struct {
	repo IEventRepository
}

func NewEventService(repo IEventRepository) *EventService {
	return &EventService{repo: repo}
}

func (service *EventService) Save(event *Event) {
	service.repo.Save(event)
}

func (service *EventService) Read(sagaId uuid.UUID, action string) (*Event, error) {
	return service.repo.Read(sagaId, action)
}
func (service *EventService) Delete(sagaId uuid.UUID, action string) error {
	return service.repo.Delete(sagaId, action)

}

func (service *EventService) GetEventEntity(event *Event, entity interface{}) {
	userProfileBytes, _ := bson.Marshal(event.Entity)
	bson.Unmarshal(userProfileBytes, entity)
}
