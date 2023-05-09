package service

import (
	"accommodation_service/domain/model"
	"accommodation_service/domain/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AccommodationService struct {
	accommodationRepo repository.IAccommodationRepository
}

func NewAccommodationService(accommodationRepo repository.IAccommodationRepository) *AccommodationService {
	return &AccommodationService{accommodationRepo: accommodationRepo}
}

func (service AccommodationService) Create(accommodation *model.Accommodation) (primitive.ObjectID, error) {
	return service.accommodationRepo.Create(accommodation)
}

func (service AccommodationService) GetById(id primitive.ObjectID) (*model.Accommodation, error) {
	return service.accommodationRepo.GetById(id)
}

func (service AccommodationService) Update(id primitive.ObjectID, dto *model.Accommodation) (*model.Accommodation, error) {
	accommodation, err := service.GetById(id)
	if err != nil {
		return &model.Accommodation{}, err
	}

	accommodation.Name = dto.Name

	accommodation, err = service.accommodationRepo.Update(accommodation)
	if err != nil {
		return nil, err
	}

	return &model.Accommodation{
		Name: accommodation.Name,
	}, nil
}
func (service AccommodationService) Delete(id primitive.ObjectID) error {
	return service.accommodationRepo.Delete(id)
}

func (service AccommodationService) GetAll() (model.Accommodations, error) {
	return service.accommodationRepo.GetAll()
}
