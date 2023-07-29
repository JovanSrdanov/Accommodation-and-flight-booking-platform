package startup

import (
	"accommodation_service/communication/handler"
	"accommodation_service/communication/middleware"
	"accommodation_service/domain/service"
	"accommodation_service/persistence/repository"
	"common/event_sourcing"
	accommodation "common/proto/accommodation_service/generated"
	"common/saga/messaging"
	"common/saga/messaging/nats"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"google.golang.org/grpc/credentials"
	"io/ioutil"
	"log"
	"net"

	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
)

const (
	QueueGroup = "accommodation_service"
)

type Server struct {
	config *Configuration
}

func NewServer(config *Configuration) *Server {
	return &Server{config: config}
}

func (server *Server) Start() {
	mongoClient := server.initMongoClient()
	accommodationRepo := initUserProfileRepo(mongoClient)
	accommodationService := service.NewAccommodationService(*accommodationRepo)

	reservationServiceAddress := fmt.Sprintf("%s:%s", server.config.ReservationServiceHost, server.config.ReservationServicePort)

	accommodationHandler := handler.NewAccommodationHandler(*accommodationService, reservationServiceAddress)

	commandSubscriber := server.initDeleteSubscriber(server.config.DeleteUserCommandSubject, QueueGroup)
	replyPublisher := server.initDeletePublisher(server.config.DeleteUserReplySubject)

	mongoEventStoreClient := server.initEventStoreMongoClient()
	eventRepo := server.initEventRepo(mongoEventStoreClient)
	eventService := event_sourcing.NewEventService(eventRepo)

	server.initDeleteHandler(accommodationService, replyPublisher, commandSubscriber, eventService)

	server.startGrpcServer(accommodationHandler)
}

func (server *Server) initEventRepo(client *mongo.Client) *event_sourcing.EventRepositoryMongo {
	repo, err := event_sourcing.NewEventRepositoryMongo(client, server.config.AccommodationEventInnerDbName, server.config.AccommodationEventDbCollectionName)
	if err != nil {
		log.Fatal(err)
	}
	return repo
}
func (server *Server) initMongoClient() *mongo.Client {
	client, err := repository.GetClient(server.config.DBName, server.config.DBPort)
	if err != nil {
		log.Fatal(err)
	}
	return client
}
func (server *Server) initEventStoreMongoClient() *mongo.Client {
	client, err := repository.GetClient(server.config.AccommodationEventDbName, server.config.AccommodationEventDbPort)
	if err != nil {
		log.Fatal(err)
	}
	return client
}

func initUserProfileRepo(mongoClient *mongo.Client) *repository.AccommodationRepositoryMongo {
	repo, err := repository.NewAccommodationRepositoryMongo(mongoClient)
	if err != nil {
		log.Fatal(err)
	}
	return repo
}

func loadTLSCredentials() (credentials.TransportCredentials, error) {
	// Load server's certificate and private key
	serverCert, err := tls.LoadX509KeyPair("/root/cert/accommodation-service-cert.pem", "/root/cert/accommodation-service-key.pem")
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

func (server *Server) startGrpcServer(userProfileHandler *handler.AccommodationHandler) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", server.config.Port))
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
		grpc.ChainUnaryInterceptor(middleware.NewGRPUnaryServerInterceptor() /*, authInterceptor.Unary()*/),
		/*grpc.StreamInterceptor(authInterceptor.Stream()),*/
	)

	accommodation.RegisterAccommodationServiceServer(grpcServer, userProfileHandler)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}

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

func (server *Server) initDeleteHandler(service *service.AccommodationService, publisher messaging.Publisher, subscriber messaging.Subscriber, eventService *event_sourcing.EventService) {
	_, err := handler.NewDeleteAccomodationHandler(service, publisher, subscriber, eventService)
	if err != nil {
		log.Fatal(err)
	}
}

// returns a map which consists of a list of grpc methods and allowed roles for each of them
//func getProtectedMethodsWithAllowedRoles() map[string][]model.Role {
//	const authServicePath = "/accommodation.AccommodationService/"
//
//	return map[string][]model.Role{
//		authServicePath + "GetAllMy":       {model.Host},
//		authServicePath + "Create":         {model.Host},
//		authServicePath + "DeleteByHostId": {model.Host},
//	}
//}
