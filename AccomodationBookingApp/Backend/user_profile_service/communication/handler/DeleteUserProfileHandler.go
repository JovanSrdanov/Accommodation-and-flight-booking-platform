package handler

import (
	accommodation "common/proto/accommodation_service/generated"
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
	userProfileService          *service.UserProfileService
	reservationServiceAddress   string
	accommodationServiceAddress string
	eventService                *event_sourcing.EventService
	replyPublisher              messaging.Publisher
	commandSubscriber           messaging.Subscriber
}

func NewDeleteUserProfileHandler(userProfileService *service.UserProfileService, reservationServiceAddress, accommodationServiceAddress string, eventService *event_sourcing.EventService, replyPublisher messaging.Publisher, commandSubscriber messaging.Subscriber) (*DeleteUserProfileHandler, error) {
	handler := &DeleteUserProfileHandler{
		userProfileService:          userProfileService,
		reservationServiceAddress:   reservationServiceAddress,
		accommodationServiceAddress: accommodationServiceAddress,
		eventService:                eventService,
		replyPublisher:              replyPublisher,
		commandSubscriber:           commandSubscriber}

	err := handler.commandSubscriber.Subscribe(handler.handle)
	if err != nil {
		return nil, err
	}

	return handler, nil
}

func (handler *DeleteUserProfileHandler) GuestHasActiveReservations(command *events.DeleteUserCommand) (bool, error) {
	reservationClient := communication.NewReservationClient(handler.reservationServiceAddress)
	result, err := reservationClient.GuestHasActiveReservations(context.TODO(), &reservation.GuestHasActiveReservationsRequest{GuestId: command.AccCredId})
	if err != nil {
		return false, err
	}
	return result.HasActiveReservations, nil
}

func (handler *DeleteUserProfileHandler) HostHasActiveReservations(command *events.DeleteUserCommand) (bool, error) {
	reservationClient := communication.NewReservationClient(handler.reservationServiceAddress)
	result, err := reservationClient.HostHasActiveReservations(context.TODO(), &reservation.HostHasActiveReservationsRequest{HostId: command.AccCredId})
	if err != nil {
		return false, err
	}

	return result.HasActiveReservations, nil
}

func (handler *DeleteUserProfileHandler) DeleteUserProfile(command *events.DeleteUserCommand) error {
	//For rollback if needed
	userProfileBackup, err := handler.userProfileService.GetById(command.UserProfileId)
	if err != nil {
		return err
	}

	err = handler.userProfileService.Delete(command.UserProfileId)
	if err != nil {
		return err
	}

	//For rollback if needed
	handler.eventService.Save(&event_sourcing.Event{
		ID:     primitive.NewObjectID(),
		SagaId: command.SagaId,
		Action: command.Type.String(),
		Entity: userProfileBackup,
	})

	return nil
}

func (handler *DeleteUserProfileHandler) RollbackProfile(deleteAction string, command *events.DeleteUserCommand) {
	//TODO error handling
	deleteEvent, _ := handler.eventService.Read(command.SagaId, deleteAction)

	var userProfile model.UserProfile
	handler.eventService.GetEventEntity(deleteEvent, &userProfile)
	handler.userProfileService.Create(&userProfile)
}

func (handler *DeleteUserProfileHandler) handle(command *events.DeleteUserCommand) {
	reply := events.DeleteUserReply{
		SagaId:        command.SagaId,
		AccCredId:     command.AccCredId,
		UserProfileId: command.UserProfileId,
		ErrorMessage:  "",
		Type:          events.UnknownReply,
	}

	switch command.Type {
	case events.DeleteGuestProfile:
		guestHasActiveReservations, err := handler.GuestHasActiveReservations(command)
		if err != nil {
			reply.Type = events.GuestProfileDeletionFailed
			reply.ErrorMessage = err.Error()
			break
		}
		if guestHasActiveReservations {
			reply.Type = events.GuestProfileDeletionFailed
			reply.ErrorMessage = "Guest has active reservations"
			break
		}

		err = handler.DeleteUserProfile(command)
		if err != nil {
			reply.Type = events.GuestProfileDeletionFailed
			break
		}

		reply.Type = events.DeletedGuestProfile
		break
	case events.DeleteHostProfile:
		hostHasActiveReservations, err := handler.HostHasActiveReservations(command)
		if err != nil {
			reply.Type = events.HostProfileDeletionFailed
			break
		}

		if hostHasActiveReservations {
			reply.Type = events.HostProfileDeletionFailed
			reply.ErrorMessage = "Host has active reservations"
			break
		}

		err = handler.DeleteUserProfile(command)
		if err != nil {
			reply.Type = events.HostProfileDeletionFailed
			break
		}

		//Delete accommodations
		accommodationClient := communication.NewAccommodationClient(handler.accommodationServiceAddress)
		result, err := accommodationClient.DeleteAllByHostId(context.TODO(), &accommodation.DeleteAllByHostIdRequest{HostId: command.AccCredId})
		if err != nil {
			reply.Type = events.HostProfileDeletionFailed
			break
		}

		accommodationIds := result.AccommodationIds

		//Delete reservations and availabilities
		reservationClient := communication.NewReservationClient(handler.reservationServiceAddress)
		accommodationServiceResult, err := reservationClient.DeleteAvailabilitiesAndReservationsByAccommodationIds(context.TODO(), &reservation.DeleteAvailabilitiesAndReservationsByAccommodationIdsRequest{AccommodationIds: accommodationIds})
		if err != nil || !accommodationServiceResult.Success {
			reply.Type = events.HostProfileDeletionFailed
			break
		}

		reply.Type = events.DeletedHostProfile
		break
	case events.RollbackGuestProfile:
		handler.RollbackProfile(events.DeleteGuestProfile.String(), command)
		reply.Type = events.RolledbackGuestProfile
		break
	case events.RollbackHostProfile:
		handler.RollbackProfile(events.DeleteHostProfile.String(), command)
		reply.Type = events.RolledbackHostProfile
		break
	default:
		reply.Type = events.UnknownReply
	}

	if reply.Type != events.UnknownReply {
		_ = handler.replyPublisher.Publish(reply)
	}
}
