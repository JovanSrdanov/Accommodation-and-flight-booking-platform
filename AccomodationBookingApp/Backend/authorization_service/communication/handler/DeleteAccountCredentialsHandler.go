package handler

import (
	"authorization_service/domain/service"
	events "common/saga/delete_user"
	"common/saga/messaging"
)

type DeleteAccountCredentialsHandler struct {
	accountCredentialsService *service.AccountCredentialsService
	replyPublisher            messaging.Publisher
	commandSubscriber         messaging.Subscriber
}

func NewDeleteAccountCredentialsHandler(accountCredentialsService *service.AccountCredentialsService, replyPublisher messaging.Publisher, commandSubscriber messaging.Subscriber) (*DeleteAccountCredentialsHandler, error) {
	handler := &DeleteAccountCredentialsHandler{
		accountCredentialsService: accountCredentialsService,
		replyPublisher:            replyPublisher,
		commandSubscriber:         commandSubscriber}

	err := handler.commandSubscriber.Subscribe(handler.handle)
	if err != nil {
		return nil, err
	}
	return handler, nil
}

func (handler *DeleteAccountCredentialsHandler) handle(command *events.DeleteUserCommand) {
	reply := events.DeleteUserReply{
		SagaId:        command.SagaId,
		AccCredId:     command.AccCredId,
		UserProfileId: command.UserProfileId,
		ErrorMessage:  "",
		Type:          events.UnknownReply,
	}

	switch command.Type {
	case events.DeleteGuestAccountCredentials:
		err := handler.accountCredentialsService.Delete(command.UserProfileId)
		if err != nil {
			reply.Type = events.GuestAccountCredentialsDeletionFailed
			reply.ErrorMessage = err.Error()
			break
		}

		reply.Type = events.DeletedGuestAccountCredentials
		break
	case events.DeleteHostAccountCredentials:
		err := handler.accountCredentialsService.Delete(command.UserProfileId)
		if err != nil {
			reply.Type = events.HostAccountCredentialsDeletionFailed
			reply.ErrorMessage = err.Error()
			break
		}

		reply.Type = events.DeletedHostAccountCredentials
		break
	}

	if reply.Type != events.UnknownReply {
		_ = handler.replyPublisher.Publish(reply)
	}

}
