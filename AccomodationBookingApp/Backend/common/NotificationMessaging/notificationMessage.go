package NotificationMessaging

import "github.com/google/uuid"

type NotificationMessage struct {
	MessageType            string
	MessageForNotification string
	AccountID              uuid.UUID
}
