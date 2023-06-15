package handler

import (
	"accommodation_service/domain/model"
	"accommodation_service/domain/service"
	"common/event_sourcing"
	events "common/saga/delete_user"
	"common/saga/messaging"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DeleteAccomodationHandler struct {
	accommodationService *service.AccommodationService
	replyPublisher       messaging.Publisher
	commandSubscriber    messaging.Subscriber
	eventService         *event_sourcing.EventService
}

func NewDeleteAccomodationHandler(
	accomodationService *service.AccommodationService,
	replyPublisher messaging.Publisher,
	commandSubscriber messaging.Subscriber,
	eventService *event_sourcing.EventService) (*DeleteAccomodationHandler, error) {
	handler := &DeleteAccomodationHandler{
		accommodationService: accomodationService,
		replyPublisher:       replyPublisher,
		commandSubscriber:    commandSubscriber,
		eventService:         eventService,
	}

	err := handler.commandSubscriber.Subscribe(handler.handle)
	if err != nil {
		return nil, err
	}

	return handler, nil
}
func (handler *DeleteAccomodationHandler) handle(command *events.DeleteUserCommand) {
	reply := events.DeleteUserReply{
		SagaId:         command.SagaId,
		AccCredId:      command.AccCredId,
		UserProfileId:  command.UserProfileId,
		Response:       command.LastResponse,
		AdditionalData: nil,
		Type:           events.UnknownReply,
	}

	switch command.Type {
	case events.DeleteHostAccommodations:
		//For rollback
		backup, err := handler.accommodationService.GetAllMy(command.AccCredId)
		if backup == nil {
			reply.Type = events.DeletedHostAccommodations
			break
		}
		if err != nil {
			reply.Response.ErrorHappened = true
			reply.Response.Message = err.Error()
			reply.Type = events.HostAccommodationsDeletionFailed
			break
		}

		var accommodationIds []string
		accommodationIds, err = handler.accommodationService.DeleteAllByHostId(command.AccCredId)

		if err != nil {
			reply.Response.ErrorHappened = true
			reply.Response.Message = err.Error()
			reply.Type = events.HostAccommodationsDeletionFailed
			break
		}

		//For rollback if needed
		handler.eventService.Save(&event_sourcing.Event{
			ID:     primitive.NewObjectID(),
			SagaId: command.SagaId,
			Action: command.Type.String(),
			Entity: backup,
		})

		reply.AdditionalData = accommodationIds
		reply.Type = events.DeletedHostAccommodations
		break
	case events.RollbackHostAccommodations:
		deleteAction := events.DeleteHostAccommodations.String()
		deleteEvent, _ := handler.eventService.Read(command.SagaId, deleteAction)

		//Konvertovanje liste entiteta iz eventa u odgovarajuce entitete
		accommodations := make([]model.Accommodation, 0)
		for _, entity := range deleteEvent.Entity.(primitive.A) {
			var accommodation model.Accommodation
			handler.eventService.ResolveEventEntity(entity, &accommodation)
			accommodations = append(accommodations, accommodation)
		}

		for _, accommodation := range accommodations {
			handler.accommodationService.Create(&accommodation)
		}

		reply.Type = events.RolledbackHostAccommodations
		break
	default:
		reply.Type = events.UnknownReply
	}

	if reply.Type != events.UnknownReply {
		_ = handler.replyPublisher.Publish(reply)
	}
}
