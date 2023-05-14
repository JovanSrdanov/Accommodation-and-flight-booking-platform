package handler

import (
	reservation "common/proto/reservation_service/generated"
	events "common/saga/delete_user"
	"common/saga/messaging"
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"user_profile_service/communication"
	"user_profile_service/domain/model"
	"user_profile_service/domain/service"
	"user_profile_service/event_sourcing"
)

type DeleteUserProfileHandler struct {
	userProfileService        *service.UserProfileService
	reservationServiceAddress string
	eventService              *event_sourcing.EventService
	replyPublisher            messaging.Publisher
	commandSubscriber         messaging.Subscriber
}

func NewDeleteUserProfileHandler(userProfileService *service.UserProfileService, reservationServiceAddress string, eventService *event_sourcing.EventService, replyPublisher messaging.Publisher, commandSubscriber messaging.Subscriber) (*DeleteUserProfileHandler, error) {
	handler := &DeleteUserProfileHandler{
		userProfileService:        userProfileService,
		reservationServiceAddress: reservationServiceAddress,
		eventService:              eventService,
		replyPublisher:            replyPublisher,
		commandSubscriber:         commandSubscriber}

	err := handler.commandSubscriber.Subscribe(handler.handle)
	if err != nil {
		return nil, err
	}

	return handler, nil
}

func (handler *DeleteUserProfileHandler) handle(command *events.DeleteUserCommand) {
	reply := events.DeleteUserReply{
		SagaId:        command.SagaId,
		UserProfileId: command.UserProfileId,
		ErrorMessage:  "",
		Type:          events.UnknownReply,
	}

	switch command.Type {
	case events.DeleteGuestProfile:
		reservationClient := communication.NewReservationClient(handler.reservationServiceAddress)
		//TODO context namesti
		result, err := reservationClient.GuestHasActiveReservations(context.TODO(), &reservation.GuestHasActiveReservationsRequest{GuestId: command.UserProfileId.String()})
		if err != nil {

			reply.Type = events.UserProfileDeletionFailed
			reply.ErrorMessage = err.Error()
			break
		}

		if result.HasActiveReservations {
			reply.Type = events.UserProfileDeletionFailed
			reply.ErrorMessage = "Guest has active reservations"
			break
		}

		userProfileBackup, err := handler.userProfileService.GetById(command.UserProfileId)
		if err != nil {
			reply.Type = events.UserProfileDeletionFailed
			break
		}

		err = handler.userProfileService.Delete(command.UserProfileId)
		if err != nil {
			reply.Type = events.UserProfileDeletionFailed
			break
		}

		handler.eventService.Save(&event_sourcing.Event{
			ID:     primitive.NewObjectID(),
			SagaId: command.SagaId,
			Action: command.Type.String(),
			Entity: userProfileBackup,
		})
		reply.Type = events.DeletedUserProfile

		break
	case events.RollbackUserProfile:
		//TODO error handling
		deleteEvent, _ := handler.eventService.Read(command.SagaId, events.DeleteGuestProfile.String())

		var userProfile model.UserProfile
		handler.eventService.GetEventEntity(deleteEvent, &userProfile)
		handler.userProfileService.Create(&userProfile)
		reply.Type = events.RolledbackUserProfile
		break
	default:
		reply.Type = events.UnknownReply
	}

	if reply.Type != events.UnknownReply {
		_ = handler.replyPublisher.Publish(reply)
	}
}
