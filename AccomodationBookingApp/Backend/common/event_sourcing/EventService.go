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

// For some reason you can't directly read interface from bson db read
func (service *EventService) ResolveEventEntity(unresolvedEntity interface{}, entity interface{}) {
	entityBytes, _ := bson.Marshal(unresolvedEntity)
	bson.Unmarshal(entityBytes, entity)
}
