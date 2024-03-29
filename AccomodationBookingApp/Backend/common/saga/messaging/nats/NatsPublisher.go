package nats

import "github.com/nats-io/nats.go"

type Publisher struct {
	conn    *nats.EncodedConn
	subject string
}

func NewNATSPublisher(host, port, user, password, subject string) (*Publisher, error) {
	conn, err := getConnection(host, port, user, password)
	encConn, err := nats.NewEncodedConn(conn, nats.JSON_ENCODER)
	if err != nil {
		return nil, err
	}
	return &Publisher{
		conn:    encConn,
		subject: subject,
	}, nil
}

func (publisher *Publisher) Publish(message interface{}) error {
	err := publisher.conn.Publish(publisher.subject, message)
	if err != nil {
		return err
	}
	return nil
}
