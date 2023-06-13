// Code generated by "stringer -type=DeleteUserCommandType"; DO NOT EDIT.

package delete_user

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[DeleteGuestProfile-0]
	_ = x[DeleteHostProfile-1]
	_ = x[RollbackGuestProfile-2]
	_ = x[RollbackHostProfile-3]
	_ = x[DeleteGuestAccountCredentials-4]
	_ = x[DeleteHostAccountCredentials-5]
	_ = x[FinishDeletion-6]
	_ = x[UnknownCommand-7]
}

const _DeleteUserCommandType_name = "DeleteGuestProfileDeleteHostProfileRollbackGuestProfileRollbackHostProfileDeleteGuestAccountCredentialsDeleteHostAccountCredentialsFinishDeletionUnknownCommand"

var _DeleteUserCommandType_index = [...]uint8{0, 18, 35, 55, 74, 103, 131, 145, 159}

func (i DeleteUserCommandType) String() string {
	if i < 0 || i >= DeleteUserCommandType(len(_DeleteUserCommandType_index)-1) {
		return "DeleteUserCommandType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _DeleteUserCommandType_name[_DeleteUserCommandType_index[i]:_DeleteUserCommandType_index[i+1]]
}
