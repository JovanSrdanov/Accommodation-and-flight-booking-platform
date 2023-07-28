package startup

import (
	"authorization_service/communication/handler"
	"authorization_service/communication/middleware"
	"authorization_service/domain/service"
	"authorization_service/domain/token"
	"authorization_service/persistence/repository"
	"common/event_sourcing"
	authorization "common/proto/authorization_service/generated"
	"common/saga/messaging"
	"common/saga/messaging/nats"
	"crypto/tls"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"gorm.io/gorm"
	"log"
	"net"
)

const (
	QueueGroup = "authorization_service"
)

type Server struct {
	config *Configuration
}

func NewServer(config *Configuration) *Server {
	return &Server{config: config}
}

func (server *Server) Start() {
	postgresClient := server.initPostgresClient()
	accountCredentialsRepo := server.initAccountCredentialsRepo(postgresClient)

	// TODO Stefan: currently not working with .env file
	tokenMaker, err := token.NewPasetoMaker("12345678901234567890123456789012")
	if err != nil {
		log.Fatalf("cannot create token maker: %v", err)
	}

	accountCredentialsService := service.NewAccountCredentialsService(accountCredentialsRepo, tokenMaker)
	accountCredentialsHandler := handler.NewAccountCredentialsHandler(accountCredentialsService)

	commandSubscriber := server.initDeleteSubscriber(server.config.DeleteUserCommandSubject, QueueGroup)
	replyPublisher := server.initDeletePublisher(server.config.DeleteUserReplySubject)

	mongoClient := server.initMongoClient()
	eventRepo := server.initEventRepo(mongoClient)
	eventService := event_sourcing.NewEventService(eventRepo)

	server.initDeleteHandler(accountCredentialsService, replyPublisher, commandSubscriber, eventService)

	server.startGrpcServer(accountCredentialsHandler, tokenMaker)
}

func (server *Server) initMongoClient() *mongo.Client {
	client, err := repository.GetMongoClient(server.config.AuthorizationEventDbName, server.config.AuthorizationEventDbPort)
	if err != nil {
		log.Fatal(err)
	}
	return client
}

func (server *Server) initEventRepo(client *mongo.Client) *event_sourcing.EventRepositoryMongo {
	repo, err := event_sourcing.NewEventRepositoryMongo(client, server.config.AuthorizationEventInnerDbName, server.config.AuthorizationEventDbCollectionName)
	if err != nil {
		log.Fatal(err)
	}
	return repo
}

func (server *Server) initPostgresClient() *gorm.DB {
	client, err := repository.GetPostgresClient(
		server.config.DBHost, server.config.DBUser,
		server.config.DBPass, server.config.DBName,
		server.config.DBPort)
	if err != nil {
		log.Fatal(err)
	}
	return client
}

func (server *Server) initAccountCredentialsRepo(postgresClient *gorm.DB) *repository.AccountCredentialsRepositoryPG {
	repo, err := repository.NewAccountCredentialsRepositoryPG(postgresClient)
	if err != nil {
		log.Fatal(err)
	}
	return repo
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

func (server *Server) initDeleteHandler(service *service.AccountCredentialsService, publisher messaging.Publisher, subscriber messaging.Subscriber, eventService *event_sourcing.EventService) {
	_, err := handler.NewDeleteAccountCredentialsHandler(service, publisher, subscriber, eventService)
	if err != nil {
		log.Fatal(err)
	}
}

func loadTLSCredentials() (credentials.TransportCredentials, error) {
	// Load server's certificate and private key
	serverCert, err := tls.LoadX509KeyPair("/root/cert/auth-service-cert.pem", "/root/cert/auth-service-key.pem")
	if err != nil {
		return nil, err
	}

	// Create the credentials and return it
	config := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.NoClientCert,
	}

	return credentials.NewTLS(config), nil
}

func (server *Server) startGrpcServer(
	accountCredentialsHandler *handler.AccountCredentialsHandler,
	maker token.Maker,
) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", server.config.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	// interceptor initialization for auth
	//protectedMethodsWithAllowedRoles := getProtectedMethodsWithAllowedRoles()
	//authInterceptor := interceptor.NewAuthServerInterceptor(maker, protectedMethodsWithAllowedRoles)

	// Enable TLS
	tlsCredentials, err := loadTLSCredentials()
	if err != nil {
		log.Fatal("cannot load TLS credentials: ", err)
	}

	grpcServer := grpc.NewServer(
		grpc.Creds(tlsCredentials),
		grpc.ChainUnaryInterceptor(middleware.NewGRPUnaryServerInterceptor()), /*authInterceptor.Unary()),
		grpc.StreamInterceptor(authInterceptor.Stream()*/
	)
	authorization.RegisterAuthorizationServiceServer(grpcServer, accountCredentialsHandler)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}

// returns a map which consists of a list of grpc methods and allowed roles for each of them
//func getProtectedMethodsWithAllowedRoles() map[string][]model.Role {
//	const authServicePath = "/authorization.AuthorizationService/"
//
//	return map[string][]model.Role{
//		authServicePath + "GetByUsername":  {model.Guest, model.Host},
//		authServicePath + "ChangeUsername": {model.Guest, model.Host},
//		authServicePath + "ChangePassword": {model.Guest, model.Host},
//		authServicePath + "CheckIfDeleted": {model.Guest, model.Host},
//	}
//}
