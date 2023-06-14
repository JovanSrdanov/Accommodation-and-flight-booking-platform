package handler

import (
	"authorization_service/domain/model"
	"authorization_service/domain/service"
	"common/event_sourcing"
	events "common/saga/delete_user"
	"common/saga/messaging"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DeleteAccountCredentialsHandler struct {
	accountCredentialsService *service.AccountCredentialsService
	replyPublisher            messaging.Publisher
	commandSubscriber         messaging.Subscriber
	eventService              *event_sourcing.EventService
}

func NewDeleteAccountCredentialsHandler(accountCredentialsService *service.AccountCredentialsService, replyPublisher messaging.Publisher, commandSubscriber messaging.Subscriber, eventService *event_sourcing.EventService) (*DeleteAccountCredentialsHandler, error) {
	handler := &DeleteAccountCredentialsHandler{
		accountCredentialsService: accountCredentialsService,
		replyPublisher:            replyPublisher,
		commandSubscriber:         commandSubscriber,
		eventService:              eventService,
	}

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
		Response:      command.LastResponse,
		Type:          events.UnknownReply,
	}

	switch command.Type {
	case events.DeleteGuestAccountCredentials:

		//For rollback
		id, _ := uuid.Parse(command.AccCredId)
		backup, _ := handler.accountCredentialsService.GetById(id)

		err := handler.accountCredentialsService.Delete(command.UserProfileId)
		if err != nil {
			reply.Type = events.GuestAccountCredentialsDeletionFailed
			reply.Response.ErrorHappened = true
			reply.Response.Message = err.Error()
			break
		}

		//For rollback
		handler.eventService.Save(&event_sourcing.Event{
			ID:     primitive.NewObjectID(),
			SagaId: command.SagaId,
			Action: command.Type.String(),
			Entity: backup,
		})

		reply.Type = events.DeletedGuestAccountCredentials
		reply.Response.Message = "Deleted guest account credentials"
		break
	case events.DeleteHostAccountCredentials:
		//For rollback
		id, _ := uuid.Parse(command.AccCredId)
		backup, _ := handler.accountCredentialsService.GetById(id)

		err := handler.accountCredentialsService.Delete(command.UserProfileId)
		if err != nil {
			reply.Type = events.HostAccountCredentialsDeletionFailed
			reply.Response.ErrorHappened = true
			reply.Response.Message = err.Error()
			break
		}

		//For rollback
		handler.eventService.Save(&event_sourcing.Event{
			ID:     primitive.NewObjectID(),
			SagaId: command.SagaId,
			Action: command.Type.String(),
			Entity: backup,
		})

		reply.Type = events.DeletedHostAccountCredentials
		reply.Response.Message = "Deleted host account credentials"
		break
	case events.RollbackGuestAccountCredentials:
		//TODO error handling
		deleteAction := events.DeleteGuestAccountCredentials.String()
		deleteEvent, _ := handler.eventService.Read(command.SagaId, deleteAction)

		var guestAccountCredentials model.AccountCredentials
		handler.eventService.ResolveEventEntity(deleteEvent.Entity, &guestAccountCredentials)
		handler.accountCredentialsService.Create(&guestAccountCredentials)

		reply.Type = events.RolledbackGuestAccountCredentials
		break
	case events.RollbackHostAccountCredentials:
		//TODO error handling
		deleteAction := events.DeleteHostAccountCredentials.String()
		deleteEvent, _ := handler.eventService.Read(command.SagaId, deleteAction)

		var hostAccountCredentials model.AccountCredentials
		handler.eventService.ResolveEventEntity(deleteEvent.Entity, &hostAccountCredentials)
		handler.accountCredentialsService.Create(&hostAccountCredentials)

		reply.Type = events.RolledbackHostAccountCredentials
		break
	default:
		reply.Type = events.UnknownReply
	}

	if reply.Type != events.UnknownReply {
		_ = handler.replyPublisher.Publish(reply)
	}

}
