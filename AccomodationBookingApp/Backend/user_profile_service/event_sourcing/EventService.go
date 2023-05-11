package event_sourcing

type EventService struct {
	repo IEventRepository
}

func NewEventService(repo IEventRepository) *EventService {
	return &EventService{repo: repo}
}

func (service *EventService) Save(event *Event) {
	service.repo.Save(event)
}
