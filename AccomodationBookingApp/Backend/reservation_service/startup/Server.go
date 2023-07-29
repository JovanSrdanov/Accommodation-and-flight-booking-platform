package startup

import (
	reservation "common/proto/reservation_service/generated"
	"common/saga/messaging"
	"common/saga/messaging/nats"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"google.golang.org/grpc/credentials"
	"io/ioutil"
	"log"
	"net"
	"reservation_service/communication/handler"
	"reservation_service/communication/middleware"
	"reservation_service/domain/service"
	"reservation_service/persistence/repository"

	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
)

const (
	QueueGroup = "reservation_service"
)

type Server struct {
	config *Configuration
}

func NewServer(config *Configuration) *Server {
	return &Server{config: config}
}

func (server *Server) Start() {
	mongoClient := server.initMongoClient()

	sendEventPublisher := server.initSendEventPublisher(server.config.SendEventToNotificationServiceSubject)
	reservationRepo := initReservationRepo(mongoClient, sendEventPublisher)
	reservationService := service.NewReservationService(*reservationRepo)
	reservationHandler := handler.NewReservationHandler(*reservationService)

	commandSubscriber := server.initDeleteSubscriber(server.config.DeleteUserCommandSubject, QueueGroup)
	replyPublisher := server.initDeletePublisher(server.config.DeleteUserReplySubject)

	server.initDeleteHandler(reservationService, replyPublisher, commandSubscriber)

	server.startGrpcServer(reservationHandler)
}

func (server *Server) initMongoClient() *mongo.Client {
	client, err := repository.GetClient(server.config.DBPort, server.config.DBName)
	if err != nil {
		log.Fatal(err)
	}
	return client
}

func initReservationRepo(mongoClient *mongo.Client, publisher messaging.Publisher) *repository.ReservationRepositoryMongo {
	repo, err := repository.NewReservationRepositoryMongo(mongoClient, publisher)
	if err != nil {
		log.Fatal(err)
	}
	return repo
}

func loadTLSCredentials() (credentials.TransportCredentials, error) {
	// Load server's certificate and private key
	serverCert, err := tls.LoadX509KeyPair("/root/cert/reservation-service-cert.pem", "/root/cert/reservation-service-key.pem")
	if err != nil {
		return nil, err
	}

	// Load certificate of the CA who signed the certificate
	pemServerCA, err := ioutil.ReadFile("/root/cert/ca-cert.pem")
	if err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pemServerCA) {
		return nil, fmt.Errorf("failed to add the CA certificate")
	}

	// Create the credentials and return it
	config := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    certPool,
	}

	return credentials.NewTLS(config), nil
}

func (server *Server) startGrpcServer(reservationHandler *handler.ReservationHandler) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", server.config.Port))
	log.Println("port: " + server.config.Port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	//tokenMaker, _ := token.NewPasetoMaker("12345678901234567890123456789012")
	//protectedMethodsWithAllowedRoles := getProtectedMethodsWithAllowedRoles()
	//authInterceptor := interceptor.NewAuthServerInterceptor(tokenMaker, protectedMethodsWithAllowedRoles)

	// Enable TLS
	tlsCredentials, err := loadTLSCredentials()
	if err != nil {
		log.Fatal("cannot load TLS credentials: ", err)
	}

	grpcServer := grpc.NewServer(
		grpc.Creds(tlsCredentials),
		grpc.ChainUnaryInterceptor(middleware.NewGRPUnaryServerInterceptor())) /*authInterceptor.Unary()),
	grpc.StreamInterceptor(authInterceptor.Stream())*/

	reservation.RegisterReservationServiceServer(grpcServer, reservationHandler)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}

// // returns a map which consists of a list of grpc methods and allowed roles for each of them
//
//	func getProtectedMethodsWithAllowedRoles() map[string][]model.Role {
//		const authServicePath = "/reservation.ReservationService/"
//
//		return map[string][]model.Role{
//			authServicePath + "GetAllMy":                   {model.Host},
//			authServicePath + "GetAllPendingReservations":  {model.Host},
//			authServicePath + "GetAllAcceptedReservations": {model.Host},
//			authServicePath + "AcceptReservation":          {model.Host},
//			authServicePath + "RejectReservation":          {model.Host},
//			authServicePath + "CancelReservation":          {model.Guest},
//			authServicePath + "GetAllReservationsForGuest": {model.Guest},
//			authServicePath + "CreateReservation":          {model.Guest},
//		}
//	}
func (server *Server) initDeleteSubscriber(subject, queueGroup string) messaging.Subscriber {
	subscriber, err := nats.NewNATSSubscriber(
		server.config.NatsHost, server.config.NatsPort,
		server.config.NatsUser, server.config.NatsPass, subject, queueGroup)
	if err != nil {
		log.Fatal(err)
	}
	return subscriber
}

func (server *Server) initDeletePublisher(subject string) messaging.Publisher {
	publisher, err := nats.NewNATSPublisher(
		server.config.NatsHost, server.config.NatsPort,
		server.config.NatsUser, server.config.NatsPass, subject)
	if err != nil {
		log.Fatal(err)
	}
	return publisher
}

func (server *Server) initDeleteHandler(service *service.ReservationService, publisher messaging.Publisher,
	subscriber messaging.Subscriber) {
	_, err := handler.NewDeleteAccomodationHandler(service, publisher, subscriber)
	if err != nil {
		log.Fatal(err)
	}
}

func (server *Server) initSendEventPublisher(subject string) messaging.Publisher {
	publisher, err := nats.NewNATSPublisher(
		server.config.NatsHost, server.config.NatsPort,
		server.config.NatsUser, server.config.NatsPass, subject)
	if err != nil {
		log.Fatal(err)
	}
	return publisher
}
