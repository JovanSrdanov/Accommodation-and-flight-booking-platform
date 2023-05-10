package handler

import (
	events "common/saga/delete_user"
	"common/saga/messaging"
	"user_profile_service/domain/service"
)

type DeleteUserProfileHandler struct {
	userProfileService *service.UserProfileService
	replyPublisher     messaging.Publisher
	commandSubscriber  messaging.Subscriber
}

func NewDeleteUserProfileHandler(userProfileService *service.UserProfileService, replyPublisher messaging.Publisher, commandSubscriber messaging.Subscriber) (*DeleteUserProfileHandler, error) {
	handler := &DeleteUserProfileHandler{
		userProfileService: userProfileService,
		replyPublisher:     replyPublisher,
		commandSubscriber:  commandSubscriber}

	err := handler.commandSubscriber.Subscribe(handler.handle)
	if err != nil {
		return nil, err
	}

	return handler, nil
}

// TODO dovrsi
func (handler *DeleteUserProfileHandler) handle(command *events.DeleteUserCommand) {
	reply := events.DeleteUserReply{
		SagaId:        command.SagaId,
		UserProfileId: command.UserProfileId,
		ErrorMessage:  "",
	}

	switch command.Type {
	case events.DeleteUserProfile:
		err := handler.userProfileService.Delete(command.UserProfileId)
		//TODO sacuvati u neku rollback tabelu
		if err != nil {
			reply.Type = events.DeletedUserProfile
		}
		break
	default:
		reply.Type = events.UnknownReply
	}

	if reply.Type != events.UnknownReply {
		_ = handler.replyPublisher.Publish(reply)
	}
}
