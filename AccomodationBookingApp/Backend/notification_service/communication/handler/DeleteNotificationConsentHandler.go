package handler

import (
	"common/event_sourcing"
	events "common/saga/delete_user"
	"common/saga/messaging"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"notification_service/domain/model"
	"notification_service/domain/service"
)

type DeleteNotificationConsentHandler struct {
	notificationConsentService *service.NotificationConsentService
	eventService               *event_sourcing.EventService
	replyPublisher             messaging.Publisher
	commandSubscriber          messaging.Subscriber
}

func NewDeleteNotificationConsentHandler(notificationConsentService *service.NotificationConsentService,
	eventService *event_sourcing.EventService,
	replyPublisher messaging.Publisher,
	commandSubscriber messaging.Subscriber) (*DeleteNotificationConsentHandler, error) {
	handler := &DeleteNotificationConsentHandler{
		notificationConsentService: notificationConsentService,
		eventService:               eventService,
		replyPublisher:             replyPublisher,
		commandSubscriber:          commandSubscriber,
	}

	err := handler.commandSubscriber.Subscribe(handler.handle)
	if err != nil {
		return nil, err
	}

	return handler, nil
}

func (handler *DeleteNotificationConsentHandler) handle(command *events.DeleteUserCommand) {
	reply := events.DeleteUserReply{
		SagaId:        command.SagaId,
		AccCredId:     command.AccCredId,
		UserProfileId: command.UserProfileId,
		Response:      command.LastResponse,
		Type:          events.UnknownReply,
	}

	switch command.Type {
	case events.DeleteHostNotifications:

		id, err := uuid.Parse(command.AccCredId)
		if err != nil {
			reply.Type = events.HostNotificationsDeletionFailed
			reply.Response.ErrorHappened = true
			reply.Response.Message = err.Error()
			break
		}

		//For rollback
		backup, _ := handler.notificationConsentService.GetById(id)

		err = handler.notificationConsentService.Delete(id)
		if err != nil {
			reply.Type = events.HostNotificationsDeletionFailed
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

		reply.Type = events.DeletedHostNotifications
		break
	case events.DeleteGuestNotifications:
		id, err := uuid.Parse(command.AccCredId)
		if err != nil {
			reply.Type = events.HostNotificationsDeletionFailed
			reply.Response.ErrorHappened = true
			reply.Response.Message = err.Error()
			break
		}

		//For rollback
		backup, _ := handler.notificationConsentService.GetById(id)
		err = handler.notificationConsentService.Delete(id)

		if err != nil {
			reply.Type = events.HostNotificationsDeletionFailed
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

		reply.Type = events.DeletedGuestNotifications
		break
	case events.RollbackHostNotifications:
		deleteAction := events.DeleteHostNotifications.String()
		deleteEvent, _ := handler.eventService.Read(command.SagaId, deleteAction)

		var notificationConsent model.NotificationConsent
		handler.eventService.ResolveEventEntity(deleteEvent.Entity, &notificationConsent)
		handler.notificationConsentService.Create(&notificationConsent)

		reply.Type = events.RolledbackHostNotifications
		break
	case events.RollbackGuestNotifications:
		deleteAction := events.DeleteGuestNotifications.String()
		deleteEvent, _ := handler.eventService.Read(command.SagaId, deleteAction)

		var notificationConsent model.NotificationConsent
		handler.eventService.ResolveEventEntity(deleteEvent.Entity, &notificationConsent)
		handler.notificationConsentService.Create(&notificationConsent)

		reply.Type = events.RolledbackGuestNotifications
		break
	default:
		reply.Type = events.UnknownReply
	}

	if reply.Type != events.UnknownReply {
		_ = handler.replyPublisher.Publish(reply)
	}
}
