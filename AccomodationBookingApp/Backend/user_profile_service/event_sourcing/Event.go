package event_sourcing

import "go.mongodb.org/mongo-driver/bson/primitive"

type Event struct {
	ID     primitive.ObjectID `bson:"_id"`
	SagaId primitive.ObjectID `bson:"sagaId"`
	Action string             `bson:"action"`
	Entity interface{}        `bson:"entity"`
}
