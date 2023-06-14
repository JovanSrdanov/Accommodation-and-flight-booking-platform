package orchestrator

import (
	authorization "common/proto/authorization_service/generated"
	events "common/saga/delete_user"
	"common/saga/messaging"
	"common/saga/messaging/nats"
	"github.com/google/uuid"
	"log"
)

type NatsInfo struct {
	NatsHost string
	NatsPort string
	NatsUser string
	NatsPass string
	Subject  string
}

type DeleteUserOrchestrator struct {
	commandPublisher messaging.Publisher
	replySubscriber  messaging.Subscriber
	natsInfo         NatsInfo
}

func NewDeleteUserOrchestrator(commandPublisher messaging.Publisher, replySubscriber messaging.Subscriber, natsInfo NatsInfo) (*DeleteUserOrchestrator, error) {
	orchestrator := &DeleteUserOrchestrator{
		commandPublisher: commandPublisher,
		replySubscriber:  replySubscriber,
		natsInfo:         natsInfo,
	}

	err := orchestrator.replySubscriber.Subscribe(orchestrator.handle)

	if err != nil {
		return nil, err
	}
	return orchestrator, nil
}

func (orchestrator *DeleteUserOrchestrator) handle(reply *events.DeleteUserReply) {
	command := events.DeleteUserCommand{
		SagaId:         reply.SagaId,
		UserProfileId:  reply.UserProfileId,
		AccCredId:      reply.AccCredId,
		Type:           events.UnknownCommand,
		AdditionalData: reply.AdditionalData,
		LastResponse:   reply.Response,
	}

	command.Type = orchestrator.nextCommandType(reply.Type)
	if command.Type != events.UnknownCommand {
		err := orchestrator.commandPublisher.Publish(command)
		if err != nil {
			log.Printf(err.Error())
		}
	}
}

func (orchestrator *DeleteUserOrchestrator) Start(accCredId string, userProfileId uuid.UUID, role authorization.Role) (events.Response, error) {
	sagaId, err := uuid.NewUUID()
	if err != nil {
		return events.Response{}, err
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

	err = orchestrator.commandPublisher.Publish(command)
	if err != nil {
		return events.Response{}, err
	}

	//Waiting for saga to finish, so we can send response back

	finishChan := make(chan events.Response)

	finishSubscriber, _ := nats.NewNATSSubscriber(
		orchestrator.natsInfo.NatsHost, orchestrator.natsInfo.NatsPort,
		orchestrator.natsInfo.NatsUser, orchestrator.natsInfo.NatsPass,
		orchestrator.natsInfo.Subject, sagaId.String())
	//Every subscriber must be in different queue group, so it can handle
	//message with appropriate sagaId

	//go func() {
	err = finishSubscriber.Subscribe(GenFinishHandler(sagaId, finishChan))
	if err != nil {
		log.Println(err.Error())
	}
	//}()

	response := <-finishChan

	return response, nil
}

func GenFinishHandler(sagaId uuid.UUID, finishChan chan<- events.Response) func(command *events.DeleteUserCommand) {
	return func(command *events.DeleteUserCommand) {
		if command.Type == events.FinishDeletion && command.SagaId == sagaId {
			finishChan <- command.LastResponse
		}
	}
}

func (orchestrator *DeleteUserOrchestrator) nextCommandType(reply events.DeleteUserReplyType) events.DeleteUserCommandType {
	switch reply {
	//DELETION PATH
	case events.DeletedGuestProfile:
		return events.DeleteGuestAccountCredentials

	//case events.DeletedGuestAccountCredentials:
	//	return events.DeleteGuestNotifications
	case events.DeletedGuestAccountCredentials:
		return events.FinishDeletion

	case events.DeletedGuestNotifications:
		return events.FinishDeletion

	case events.DeletedHostProfile:
		return events.DeleteHostAccountCredentials
	case events.DeletedHostAccountCredentials:
		return events.DeleteHostAccommodations
	case events.DeletedHostAccommodations:
		return events.DeleteHostReservations

	//case events.DeletedHostReservations:
	//	return events.DeleteHostNotifications
	case events.DeletedHostReservations:
		return events.FinishDeletion

	case events.DeletedHostNotifications:
		return events.FinishDeletion
	// FAIL PATH

	case events.GuestProfileDeletionFailed:
		return events.FinishDeletion
	case events.HostProfileDeletionFailed:
		return events.FinishDeletion

	case events.GuestAccountCredentialsDeletionFailed:
		return events.RollbackGuestProfile
	case events.HostAccountCredentialsDeletionFailed:
		return events.RollbackHostProfile

	case events.HostAccommodationsDeletionFailed:
		return events.RollbackHostAccountCredentials

	case events.HostReservationsDeletionFailed:
		return events.RollbackHostAccommodations

	case events.HostNotificationsDeletionFailed:
		return events.RollbackHostReservations
	case events.GuestNotificationsDeletionFailed:
		return events.RollbackGuestAccountCredentials

	//ROLLBACK PATH
	case events.RolledbackGuestProfile:
		return events.FinishDeletion
	case events.RolledbackHostProfile:
		return events.FinishDeletion

	case events.RolledbackGuestAccountCredentials:
		return events.RollbackGuestProfile
	case events.RolledbackHostAccountCredentials:
		return events.RollbackHostProfile

	case events.RolledbackHostAccommodations:
		return events.RollbackHostAccountCredentials

	case events.RolledbackHostReservations:
		return events.RollbackHostAccommodations

	default:
		return events.UnknownCommand
	}
}
