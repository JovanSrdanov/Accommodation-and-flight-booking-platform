package event_sourcing

type IEventRepository interface {
	Save(event *Event) error
}
