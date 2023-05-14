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
		From: time.Unix(request.Reservation.DateRange.From, 0).In(time.UTC),
		To:   time.Unix(request.Reservation.DateRange.To, 0).In(time.UTC),
	}

	NormalizeTime(&dateRange)

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
		From: time.Unix(request.PriceWithDate.UpdatedPriceWithDate.DateRange.From, 0).In(time.UTC),
		To:   time.Unix(request.PriceWithDate.UpdatedPriceWithDate.DateRange.To, 0).In(time.UTC),
	}

	NormalizeTime(&dateRange)

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

func (mapper ReservationMapper) mapToReservationsProto(in model.Reservations) []*reservation.Reservation {
	reservationsProt := make([]*reservation.Reservation, 0)

	for _, reservationValue := range in {
		reservationProto := &reservation.Reservation{
			Id:     reservationValue.ID.String(),
			Status: reservationValue.Status,
			DateRange: &reservation.DateRange{
				From: reservationValue.DateRange.From.Unix(),
				To:   reservationValue.DateRange.To.Unix(),
			},
			AccommodationId: reservationValue.AccommodationId.String(),
			Price:           reservationValue.Price,
			NumberOfGuests:  reservationValue.NumberOfGuests,
			GuestId:         reservationValue.GuestId,
		}

		reservationsProt = append(reservationsProt, reservationProto)
	}

	return reservationsProt
}

func (mapper ReservationMapper) mapFromSearchRequest(request *reservation.SearchRequest) ([]*primitive.ObjectID, model.DateRange, int32) {
	ids := make([]*primitive.ObjectID, 0)

	for _, protoId := range request.Filter.AccommodationIds {
		conv, _ := primitive.ObjectIDFromHex(protoId)
		ids = append(ids, &conv)
	}

	dateRange := model.DateRange{
		From: time.Unix(request.Filter.DateRange.From, 0).In(time.UTC),
		To:   time.Unix(request.Filter.DateRange.To, 0).In(time.UTC),
	}

	NormalizeTime(&dateRange)

	return ids, dateRange, request.Filter.NumberOfGuests
}

func (mapper ReservationMapper) mapToSearchResponse(in []*model.SearchResponseDto) *reservation.SearchResponse {
	searchResponseProto := make([]*reservation.SearchResponseDto, 0)

	for _, dto := range in {
		searchResponseProto = append(searchResponseProto, &reservation.SearchResponseDto{
			AccommodationId: dto.AccommodationId.Hex(),
			Price:           dto.Price,
		})
	}

	return &reservation.SearchResponse{SearchResponse: searchResponseProto}
}

func NormalizeTime(dateRange *model.DateRange) {
	dateRange.To = dateRange.To.In(time.UTC)
	dateRange.From = dateRange.From.In(time.UTC)

	dateRange.To = time.Date(dateRange.To.Year(), dateRange.To.Month(), dateRange.To.Day(), 0, 0, 0, 0, dateRange.To.Location())
	dateRange.From = time.Date(dateRange.From.Year(), dateRange.From.Month(), dateRange.From.Day(), 0, 0, 0, 0, dateRange.From.Location())

	//dateRange.To = dateRange.To.AddDate(0, 0, 1)
	//dateRange.From = dateRange.From.AddDate(0, 0, 1)
}
