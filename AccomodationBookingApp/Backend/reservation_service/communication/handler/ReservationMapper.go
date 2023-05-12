package handler

import (
	reservation "common/proto/reservation_service/generated"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"reservation_service/domain/model"
	"time"
)

type ReservationMapper struct{}

type IReservationMapper interface {
	mapFromCreateAvailability(request *reservation.CreateAvailabilityRequest) *model.AvailabilityRequest
}

func NewReservationMapper() *ReservationMapper {
	return &ReservationMapper{}
}

func (mapper ReservationMapper) mapFromCreateAvailability(request *reservation.CreateAvailabilityRequest) *model.AvailabilityRequest {
	priceWithDateMapper := NewPriceWithDateMapper()

	accommodationId, _ := primitive.ObjectIDFromHex(request.Availability.AccommodationId)

	return &model.AvailabilityRequest{
		PriceWithDate:   priceWithDateMapper.mapToModel(request.Availability.PriceWithDate),
		AccommodationId: accommodationId,
	}
}

func (mapper ReservationMapper) mapFromCreateReservation(request *reservation.CreateReservationRequest) *model.Reservation {
	dateRange := model.DateRange{
		From: time.Unix(request.Reservation.DateRange.From, 0),
		To:   time.Unix(request.Reservation.DateRange.To, 0),
	}

	id, _ := primitive.ObjectIDFromHex(request.Reservation.Id)
	accommodationId, _ := primitive.ObjectIDFromHex(request.Reservation.AccommodationId)

	return &model.Reservation{
		ID:              id,
		AccommodationId: accommodationId,
		DateRange:       dateRange,
		Price:           request.Reservation.Price,
		NumberOfGuests:  request.Reservation.NumberOfGuests,
		Status:          request.Reservation.Status,
	}
}

func (mapper ReservationMapper) mapFromUpdatePriceAndDate(request *reservation.UpdateRequest) *model.UpdatePriceAndDate {
	dateRange := model.DateRange{
		From: time.Unix(request.PriceWithDate.UpdatedPriceWithDate.DateRange.From, 0),
		To:   time.Unix(request.PriceWithDate.UpdatedPriceWithDate.DateRange.To, 0),
	}

	id, _ := primitive.ObjectIDFromHex(request.PriceWithDate.UpdatedPriceWithDate.Id)
	accommodationId, _ := primitive.ObjectIDFromHex(request.PriceWithDate.AccommodationId)

	return &model.UpdatePriceAndDate{
		AccommodationId: accommodationId,
		PriceWithDate: model.PriceWithDate{
			ID:               id,
			DateRange:        dateRange,
			Price:            request.PriceWithDate.UpdatedPriceWithDate.Price,
			IsPricePerPerson: request.PriceWithDate.UpdatedPriceWithDate.IsPricePerPerson,
		},
	}
}

func (mapper ReservationMapper) mapFromCreateAvailabilityBase(hostId string, request *reservation.CreateAvailabilityBaseRequest) *model.Availability {
	accommodationId, _ := primitive.ObjectIDFromHex(request.ReservationBase.AccommodationId)

	return &model.Availability{
		AccommodationId:        accommodationId,
		HostId:                 hostId,
		IsAutomaticReservation: request.ReservationBase.IsAutomaticReservation,
	}
}

func (mapper ReservationMapper) mapToGetAllMyResponse(in model.Availabilities) *reservation.GetAllMyResponse {
	availabilitiesProto := make([]*reservation.Availability, 0)

	for _, avail := range in {

		priceWithDateProt := make([]*reservation.PriceWithDate, 0)

		for _, priceWithDate := range avail.AvailableDates {
			priceProto := &reservation.PriceWithDate{
				Id: priceWithDate.ID.String(),
				DateRange: &reservation.DateRange{
					From: priceWithDate.DateRange.From.Unix(),
					To:   priceWithDate.DateRange.To.Unix(),
				},
				Price:            priceWithDate.Price,
				IsPricePerPerson: priceWithDate.IsPricePerPerson,
			}

			priceWithDateProt = append(priceWithDateProt, priceProto)
		}

		availProto := &reservation.Availability{
			Id:                     avail.ID.String(),
			AvailableDates:         priceWithDateProt,
			AccommodationId:        avail.AccommodationId.String(),
			IsAutomaticReservation: avail.IsAutomaticReservation,
			HostId:                 avail.HostId,
		}
		availabilitiesProto = append(availabilitiesProto, availProto)
	}

	return &reservation.GetAllMyResponse{
		Availabilities: availabilitiesProto,
	}
}

