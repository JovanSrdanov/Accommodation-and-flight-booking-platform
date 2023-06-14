package event_sourcing

import (
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Event struct {
	ID     primitive.ObjectID `bson:"_id"`
	SagaId uuid.UUID          `bson:"sagaId"`
	Action string             `bson:"action"`
	Entity interface{}        `bson:"entity"`
}
