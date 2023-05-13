package handler

import (
	events "common/saga/delete_user"
	"common/saga/messaging"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"user_profile_service/domain/model"
	"user_profile_service/domain/service"
	"user_profile_service/event_sourcing"
)

type DeleteUserProfileHandler struct {
	userProfileService *service.UserProfileService
	eventService       *event_sourcing.EventService
	replyPublisher     messaging.Publisher
	commandSubscriber  messaging.Subscriber
}

func NewDeleteUserProfileHandler(userProfileService *service.UserProfileService, eventService *event_sourcing.EventService, replyPublisher messaging.Publisher, commandSubscriber messaging.Subscriber) (*DeleteUserProfileHandler, error) {
	handler := &DeleteUserProfileHandler{
		userProfileService: userProfileService,
		eventService:       eventService,
		replyPublisher:     replyPublisher,
		commandSubscriber:  commandSubscriber}

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
	}

	switch command.Type {
	case events.DeleteUserProfile:
		userProfileBackup, err := handler.userProfileService.GetById(command.UserProfileId)
		if err != nil {
			reply.Type = events.UserProfileDeletionFailed
		}

		err = handler.userProfileService.Delete(command.UserProfileId)
		if err != nil {
			reply.Type = events.UserProfileDeletionFailed
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
		deleteEvent, _ := handler.eventService.Read(command.SagaId, events.DeleteUserProfile.String())

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
