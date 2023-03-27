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
	BuyTicket(ticket model.Ticket, flightId primitive.ObjectID, numberOfTickets int32) (primitive.ObjectID, error)
	GetAllForCustomer() ([]dto.TicketFullInfo, error)
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

func (service *ticketService) BuyTicket(ticket model.Ticket, flightId primitive.ObjectID, numberOfTickets int32) (primitive.ObjectID, error) {
	//TODO Strahinja: buyera izvuci iz JWT
	//ownera iz api kljuca
	//Ovo treba da bude transakcija
	ticket.Buyer = "JWT"
	ticket.Owner = "APIKey"

	flight, err := service.flightRepository.GetById(flightId)
	if err != nil {
		return primitive.ObjectID{}, err
	}

	if flight.VacantSeats < numberOfTickets {
		return flightId, &errors.NotEnoughVacantSeats{}
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

	return flightId, nil
}

func (service *ticketService) GetAllForCustomer() ([]dto.TicketFullInfo, error) {
	return service.ticketRepository.GetAllForCustomer()
}
