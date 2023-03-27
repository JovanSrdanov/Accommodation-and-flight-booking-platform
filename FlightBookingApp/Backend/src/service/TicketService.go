package service

import (
	"FlightBookingApp/model"
	"FlightBookingApp/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ticketService struct {
	ticketRepository repository.TicketRepositry
}

type TicketService interface {
	Create(flight model.Ticket) (primitive.ObjectID, error)
	GetAll() (model.Tickets, error)
	GetById(id primitive.ObjectID) (model.Ticket, error)
	Delete(id primitive.ObjectID) error
}

func NewTicketService(ticketRepository repository.TicketRepositry) *ticketService {
	return &ticketService{
		ticketRepository: ticketRepository,
	}
}

func (service *ticketService) Create(ticket model.Ticket) (primitive.ObjectID, error) {
	return service.ticketRepository.Create(&ticket)
}

func (service *ticketService) GetAll() (model.Tickets, error) {
	return service.ticketRepository.GetAll()
}

func (service *ticketService) GetById(id primitive.ObjectID) (model.Ticket, error) {
	return service.ticketRepository.GetById(id)
}
func (service *ticketService) Delete(id primitive.ObjectID) error {
	return service.ticketRepository.Delete(id)
}
