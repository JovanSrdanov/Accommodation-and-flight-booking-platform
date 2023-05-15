package orchestrator

import (
	authorization "common/proto/authorization_service/generated"
	events "common/saga/delete_user"
	"common/saga/messaging"
	"github.com/google/uuid"
)

type DeleteUserOrchestrator struct {
	commandPublisher messaging.Publisher
	replySubscriber  messaging.Subscriber
}

func NewDeleteUserOrchestrator(commandPublisher messaging.Publisher, replySubscriber messaging.Subscriber) (*DeleteUserOrchestrator, error) {
	orchestrator := &DeleteUserOrchestrator{
		commandPublisher: commandPublisher,
		replySubscriber:  replySubscriber}

	err := orchestrator.replySubscriber.Subscribe(orchestrator.handle)

	if err != nil {
		return nil, err
	}
	return orchestrator, nil
}

func (orchestrator *DeleteUserOrchestrator) handle(reply *events.DeleteUserReply) {
	command := events.DeleteUserCommand{
		SagaId:        reply.SagaId,
		UserProfileId: reply.UserProfileId,
		Type:          events.UnknownCommand,
	}
	command.Type = orchestrator.nextCommandType(reply.Type)
	if command.Type != events.UnknownCommand {
		orchestrator.commandPublisher.Publish(command)
	}
}

func (orchestrator *DeleteUserOrchestrator) nextCommandType(reply events.DeleteUserReplyType) events.DeleteUserCommandType {
	switch reply {
	case events.DeletedGuestProfile:
		return events.DeleteGuestAccountCredentials
	case events.DeletedHostProfile:
		return events.DeleteHostAccountCredentials
	case events.DeletedGuestAccountCredentials:
		return events.UnknownCommand
	case events.DeletedHostAccountCredentials:
		return events.UnknownCommand
	case events.GuestProfileDeletionFailed:
		return events.UnknownCommand
	case events.HostProfileDeletionFailed:
		return events.UnknownCommand
	case events.GuestAccountCredentialsDeletionFailed:
		return events.RollbackGuestProfile
	case events.HostAccountCredentialsDeletionFailed:
		return events.RollbackHostProfile
	case events.RolledbackGuestProfile:
		return events.UnknownCommand
	case events.RolledbackHostProfile:
		return events.UnknownCommand
	default:
		return events.UnknownCommand
	}
}

func (orchestrator *DeleteUserOrchestrator) Start(accCredId string, userProfileId uuid.UUID, role authorization.Role) error {
	sagaId, err := uuid.NewUUID()
	if err != nil {
		return err
	}

	command := events.DeleteUserCommand{
		SagaId:        sagaId,
		AccCredId:     accCredId,
		UserProfileId: userProfileId,
	}

	if role == authorization.Role_Guest {
		command.Type = events.DeleteGuestProfile
	} else {
		command.Type = events.DeleteHostProfile
	}

	return orchestrator.commandPublisher.Publish(command)
}
