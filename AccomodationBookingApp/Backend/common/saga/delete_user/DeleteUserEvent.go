package delete_user

import "github.com/google/uuid"

type DeleteUserCommandType int8

const (
	DeleteUserProfile DeleteUserCommandType = iota
	RollbackUserProfile
	DeleteAccountCredentials
	CancelDeletion
	FinishDeletion
	UnknownCommand
)

type DeleteUserCommand struct {
	Type          DeleteUserCommandType
	SagaId        uuid.UUID
	UserProfileId uuid.UUID
}

type DeleteUserReplyType int8

const (
	DeletedUserProfile DeleteUserReplyType = iota
	UserProfileDeletionFailed
	DeletedAccountCredentials
	AccountCredentialsDeletionFailed
	RolledbackUserProfile
	UnknownReply
)

type DeleteUserReply struct {
	Type          DeleteUserReplyType
	SagaId        uuid.UUID
	UserProfileId uuid.UUID
	ErrorMessage  string
}
