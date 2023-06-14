package handler

import (
	events "common/saga/delete_user"
	"common/saga/messaging"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"reservation_service/domain/service"
)

type DeleteReservationHandler struct {
	reservationService *service.ReservationService
	replyPublisher     messaging.Publisher
	commandSubscriber  messaging.Subscriber
}

func NewDeleteAccomodationHandler(
	accomodationService *service.ReservationService,
	replyPublisher messaging.Publisher,
	commandSubscriber messaging.Subscriber,
) (*DeleteReservationHandler, error) {
	handler := &DeleteReservationHandler{
		reservationService: accomodationService,
		replyPublisher:     replyPublisher,
		commandSubscriber:  commandSubscriber,
	}

	err := handler.commandSubscriber.Subscribe(handler.handle)
	if err != nil {
		return nil, err
	}

	return handler, nil
}

func (handler *DeleteReservationHandler) handle(command *events.DeleteUserCommand) {
	reply := events.DeleteUserReply{
		SagaId:         command.SagaId,
		AccCredId:      command.AccCredId,
		UserProfileId:  command.UserProfileId,
		Response:       command.LastResponse,
		AdditionalData: nil,
		Type:           events.UnknownReply,
	}

	switch command.Type {
	case events.DeleteHostReservations:
		//If there is no deleted accommodations
		if command.AdditionalData == nil {
			reply.Type = events.DeletedHostReservations
			break
		}

		accomodationIdsInrterf, ok := command.AdditionalData.([]interface{})
		if !ok {
			reply.Response.ErrorHappened = true
			reply.Response.Message = "Couldn't get list of deleted accommodations"
			reply.Type = events.HostReservationsDeletionFailed
			break
		}

		accommodationIds := make([]string, 0)
		for _, idInterf := range accomodationIdsInrterf {
			accommodationIds = append(accommodationIds, idInterf.(string))
		}
		failed := false
		for _, accomodationId := range accommodationIds {
			accommodationIdObj, err := primitive.ObjectIDFromHex(accomodationId)
			if err != nil {
				reply.Response.ErrorHappened = true
				reply.Response.Message = err.Error()
				reply.Type = events.HostReservationsDeletionFailed
				failed = true
				break
			}

			err = handler.reservationService.DeleteAvailabilitiesAndReservationsByAccommodationId(accommodationIdObj)
			if err != nil {
				reply.Response.ErrorHappened = true
				reply.Response.Message = err.Error()
				reply.Type = events.HostReservationsDeletionFailed
				failed = true
				break
			}
		}
		if failed {
			break
		}

		reply.Type = events.DeletedHostReservations
		break
	default:
		reply.Type = events.UnknownReply
	}

	if reply.Type != events.UnknownReply {
		_ = handler.replyPublisher.Publish(reply)
	}
}
