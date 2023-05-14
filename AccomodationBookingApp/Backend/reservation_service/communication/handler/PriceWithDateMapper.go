package handler

import (
	reservation "common/proto/reservation_service/generated"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"reservation_service/domain/model"
	"time"
)

type PriceWithDateMapper struct{}

type IPriceWithDateMapper interface {
	mapToModel(request *reservation.PriceWithDate) *model.PriceWithDate
}

func NewPriceWithDateMapper() *PriceWithDateMapper {
	return &PriceWithDateMapper{}
}

func (mapper PriceWithDateMapper) mapToModel(request *reservation.PriceWithDate) model.PriceWithDate {
	Id, _ := primitive.ObjectIDFromHex(request.Id)
	dateRange := model.DateRange{
		From: time.Unix(request.DateRange.From, 0),
		To:   time.Unix(request.DateRange.To, 0),
	}

	NormalizeTime(&dateRange)

	return model.PriceWithDate{
		ID:               Id,
		DateRange:        dateRange,
		Price:            request.Price,
		IsPricePerPerson: request.IsPricePerPerson,
	}
}
