package service

import (
	"FlightBookingApp/dto"
	"FlightBookingApp/errors"
	"FlightBookingApp/model"
	"FlightBookingApp/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ticketService struct {
	ticketRepository repository.TicketRepositry
	flightRepository repository.FlightRepository
}

type TicketService interface {
	Create(flight model.Ticket) (primitive.ObjectID, error)
	GetAll() (model.Tickets, error)
	GetById(id primitive.ObjectID) (model.Ticket, error)
	Delete(id primitive.ObjectID) error
	BuyTicket(ticket model.Ticket, numberOfTickets int32) (primitive.ObjectID, error)
	GetAllForUser(userId primitive.ObjectID) ([]dto.TicketFullInfo, error)
}

func NewTicketService(ticketRepository repository.TicketRepositry, flightRepository repository.FlightRepository) *ticketService {
	return &ticketService{
		ticketRepository: ticketRepository,
		flightRepository: flightRepository,
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

func (service *ticketService) BuyTicket(ticket model.Ticket, numberOfTickets int32) (primitive.ObjectID, error) {
	//TODO Strahinja: ownera iz api kljuca, Ovo treba da bude transakcija

	flight, err := service.flightRepository.GetById(ticket.FlightId)
	if err != nil {
		return primitive.ObjectID{}, err
	}

	if flight.Canceled {
		return primitive.ObjectID{}, &errors.FlightIsCanceledError{}
	}

	if flight.VacantSeats < numberOfTickets {
		return primitive.ObjectID{}, &errors.NotEnoughVacantSeatsError{}
	}

	flight.DecreaseVacantSeats(numberOfTickets)

	_, err = service.flightRepository.Save(flight)
	if err != nil {
		return primitive.ObjectID{}, err
	}

	var i int32
	for i = 0; i < numberOfTickets; i++ {
		_, err = service.ticketRepository.Create(&ticket)
		if err != nil {
			return primitive.ObjectID{}, err
		}
	}

	return ticket.FlightId, nil
}

func (service *ticketService) GetAllForUser(userId primitive.ObjectID) ([]dto.TicketFullInfo, error) {
	return service.ticketRepository.GetAllForUser(userId)
}
