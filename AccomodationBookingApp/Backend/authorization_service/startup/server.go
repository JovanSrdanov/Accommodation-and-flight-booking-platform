package startup

import (
	"authorization_service/communication/handler"
	"authorization_service/domain/model"
	"authorization_service/domain/service"
	"authorization_service/domain/token"
	"authorization_service/interceptor"
	"authorization_service/persistence/repository"
	authorization "common/proto/authorization_service/generated"
	"common/saga/messaging"
	"common/saga/messaging/nats"
	"fmt"
	"google.golang.org/grpc"
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
	server.initDeleteHandler(accountCredentialsService, replyPublisher, commandSubscriber)

	server.startGrpcServer(accountCredentialsHandler, tokenMaker)
}

func (server *Server) initPostgresClient() *gorm.DB {
	client, err := repository.GetClient(
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

func (server *Server) initDeleteHandler(service *service.AccountCredentialsService, publisher messaging.Publisher, subscriber messaging.Subscriber) {
	_, err := handler.NewDeleteAccountCredentialsHandler(service, publisher, subscriber)
	if err != nil {
		log.Fatal(err)
	}
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
	protectedMethodsWithAllowedRoles := getProtectedMethodsWithAllowedRoles()
	authInterceptor := interceptor.NewAuthServerInterceptor(maker, protectedMethodsWithAllowedRoles)

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(authInterceptor.Unary()),
		grpc.StreamInterceptor(authInterceptor.Stream()),
	)
	authorization.RegisterAuthorizationServiceServer(grpcServer, accountCredentialsHandler)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}

// returns a map which consists of a list of grpc methods and allowed roles for each of them
func getProtectedMethodsWithAllowedRoles() map[string][]model.Role {
	const authServicePath = "/authorization.AuthorizationService/"

	return map[string][]model.Role{
		authServicePath + "GetByUsername":  {model.Guest, model.Host},
		authServicePath + "ChangeUsername": {model.Guest, model.Host},
		authServicePath + "ChangePassword": {model.Guest, model.Host},
	}
}
