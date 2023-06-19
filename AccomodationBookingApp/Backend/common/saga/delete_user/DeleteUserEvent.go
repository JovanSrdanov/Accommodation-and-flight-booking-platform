package delete_user

import "github.com/google/uuid"

// Da bi implementirao stringer interface, treba:https://pkg.go.dev/golang.org/x/tools/cmd/stringer
// Konkretno, u terminalu se pozicioniras u isti folder gde se enum nalazi i pokrenes:
// stringer -type=DeleteUserCommandType
type DeleteUserCommandType int8

const (
	DeleteGuestProfile DeleteUserCommandType = iota
	DeleteHostProfile
	RollbackGuestProfile
	RollbackHostProfile

	DeleteGuestAccountCredentials
	DeleteHostAccountCredentials
	RollbackGuestAccountCredentials
	RollbackHostAccountCredentials

	DeleteHostAccommodations
	RollbackHostAccommodations

	//Last in chain so don't need rollback
	DeleteHostReservations

	DeleteHostNotifications
	DeleteGuestNotifications
	RollbackHostNotifications
	RollbackGuestNotifications

	FinishDeletion
	UnknownCommand
)

type Response struct {
	ErrorHappened bool
	Message       string
}

type DeleteUserCommand struct {
	Type           DeleteUserCommandType
	SagaId         uuid.UUID
	AccCredId      string
	UserProfileId  uuid.UUID
	LastResponse   Response
	AdditionalData interface{}
}

type DeleteUserReplyType int8

const (
	DeletedGuestProfile DeleteUserReplyType = iota
	DeletedHostProfile
	GuestProfileDeletionFailed
	HostProfileDeletionFailed
	RolledbackGuestProfile
	RolledbackHostProfile

	DeletedGuestAccountCredentials
	DeletedHostAccountCredentials
	GuestAccountCredentialsDeletionFailed
	HostAccountCredentialsDeletionFailed
	RolledbackGuestAccountCredentials
	RolledbackHostAccountCredentials

	DeletedHostAccommodations
	HostAccommodationsDeletionFailed
	RolledbackHostAccommodations

	DeletedHostReservations
	HostReservationsDeletionFailed

	DeletedHostNotifications
	DeletedGuestNotifications
	HostNotificationsDeletionFailed
	GuestNotificationsDeletionFailed
	RolledbackHostNotifications
	RolledbackGuestNotifications

	UnknownReply
)

type DeleteUserReply struct {
	Type           DeleteUserReplyType
	AccCredId      string
	SagaId         uuid.UUID
	UserProfileId  uuid.UUID
	Response       Response
	AdditionalData interface{}
}
