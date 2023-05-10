package orchestrator

import (
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
	}
	command.Type = orchestrator.nextCommandType(reply.Type)
	if command.Type != events.UnknownCommand {
		orchestrator.commandPublisher.Publish(command)
	}
}

func (orchestrator *DeleteUserOrchestrator) nextCommandType(reply events.DeleteUserReplyType) events.DeleteUserCommandType {
	switch reply {
	case events.DeletedUserProfile:
		return events.DeleteAccountCredentials
	case events.DeletedAccountCredentials:
		return events.FinishDeletion
	case events.UserProfileDeletionFailed:
		return events.CancelDeletion
	case events.AccountCredentialsDeletionFailed:
		return events.RollbackUserProfile
	case events.RolledbackUserProfile:
		return events.CancelDeletion
	default:
		return events.UnknownCommand
	}
}

func (orchestrator *DeleteUserOrchestrator) Start(userProfileId uuid.UUID) error {
	sagaId, err := uuid.NewUUID()
	if err != nil {
		return err
	}

	command := events.DeleteUserCommand{
		Type:          events.DeleteUserProfile,
		SagaId:        sagaId,
		UserProfileId: userProfileId,
	}

	return orchestrator.commandPublisher.Publish(command)
}
