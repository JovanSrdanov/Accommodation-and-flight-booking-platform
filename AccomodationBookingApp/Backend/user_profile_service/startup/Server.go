package startup

import (
	"authorization_service/domain/model"
	"authorization_service/domain/token"
	"authorization_service/interceptor"
	user_profile "common/proto/user_profile_service/generated"
	"common/saga/messaging"
	"common/saga/messaging/nats"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
	"gorm.io/gorm"
	"log"
	"net"
	"user_profile_service/communication/handler"
	"user_profile_service/communication/orchestrator"
	"user_profile_service/domain/service"
	"user_profile_service/event_sourcing"
	"user_profile_service/persistence/repository"
)

const (
	QueueGroup = "user_profile_service"
)

type Server struct {
	config *Configuration
}

func NewServer(config *Configuration) *Server {
	return &Server{
		config: config,
	}
}

func (server *Server) Start() {
	postgresClient := server.initPostgresClient()
	userProfileRepo := initUserProfileRepo(postgresClient)

	//Standard handler with orchestrator
	commandPublisher := server.initDeletePublisher(server.config.DeleteUserCommandSubject)
	replySubscriber := server.initDeleteSubscriber(server.config.DeleteUserReplySubject, QueueGroup)
	deleteUserOrchestrator := server.initDeleteUserOrchestrator(commandPublisher, replySubscriber)

	userProfileService := service.NewUserProfileService(*userProfileRepo, deleteUserOrchestrator)
	authServiceAddress := fmt.Sprintf("%s:%s", server.config.AuthServiceHost, server.config.AuthServicePort)
	userProfileHandler := handler.NewUserProfileHandler(*userProfileService, authServiceAddress)

	//Delete handler that listens orchestrator
	commandSubscriber := server.initDeleteSubscriber(server.config.DeleteUserCommandSubject, QueueGroup)
	replyPublisher := server.initDeletePublisher(server.config.DeleteUserReplySubject)

	mongoClient := server.initMongoClient()
	eventRepo := server.initEventRepo(mongoClient)
	eventService := event_sourcing.NewEventService(eventRepo)

	reservationServiceAddress := fmt.Sprintf("%s:%s", server.config.ReservationServiceHost, server.config.ReservationServicePort)
	accommodationServiceAddress := fmt.Sprintf("%s:%s", server.config.AccommodationServiceHost, server.config.AccommodationServicePort)
	server.initDeleteUserHandler(userProfileService, reservationServiceAddress, accommodationServiceAddress, eventService, replyPublisher, commandSubscriber)
	server.startGrpcServer(userProfileHandler)
}

func (server *Server) initMongoClient() *mongo.Client {
	client, err := repository.GetMongoClient(server.config.UserProfileEventDbName, server.config.UserProfileEventDbPort)
	if err != nil {
		log.Fatal(err)
	}
	return client
}

func (server *Server) initEventRepo(client *mongo.Client) *repository.EventRepositoryMongo {
	repo, err := repository.NewEventRepositoryMongo(client, server.config.UserProfileEventInnerDbName, server.config.UserProfileEventDbCollectionName)
	if err != nil {
		log.Fatal(err)
	}
	return repo
}
func initUserProfileRepo(client *gorm.DB) *repository.UserProfileRepositoryPG {
	repo, err := repository.NewUserProfileRepositoryPG(client)
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
func (server *Server) startGrpcServer(userProfileHandler *handler.UserProfileHandler) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", server.config.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	tokenMaker, _ := token.NewPasetoMaker("12345678901234567890123456789012")
	protectedMethodsWithAllowedRoles := getProtectedMethodsWithAllowedRoles()
	authInterceptor := interceptor.NewAuthServerInterceptor(tokenMaker, protectedMethodsWithAllowedRoles)

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(authInterceptor.Unary()),
		grpc.StreamInterceptor(authInterceptor.Stream()),
	)
	user_profile.RegisterUserProfileServiceServer(grpcServer, userProfileHandler)
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

func (server *Server) initDeleteUserOrchestrator(publisher messaging.Publisher, subscriber messaging.Subscriber) *orchestrator.DeleteUserOrchestrator {
	orch, err := orchestrator.NewDeleteUserOrchestrator(publisher, subscriber, orchestrator.NatsInfo{
		NatsHost: server.config.NatsHost,
		NatsPort: server.config.NatsPort,
		NatsUser: server.config.NatsUser,
		NatsPass: server.config.NatsPass,
		Subject:  server.config.DeleteUserCommandSubject,
	})
	if err != nil {
		log.Fatal(err)
	}
	return orch
}

func (server *Server) initDeleteUserHandler(userProfileService *service.UserProfileService, reservationServiceAddress, accommodationServiceAddress string, eventService *event_sourcing.EventService, publisher messaging.Publisher, subscriber messaging.Subscriber) {
	_, err := handler.NewDeleteUserProfileHandler(userProfileService, reservationServiceAddress, accommodationServiceAddress, eventService, publisher, subscriber)
	if err != nil {
		log.Fatal(err)
	}

}

// returns a map which consists of a list of grpc methods and allowed roles for each of them
func getProtectedMethodsWithAllowedRoles() map[string][]model.Role {
	const authServicePath = "/user_profile.UserProfileService/"

	return map[string][]model.Role{
		authServicePath + "Update":     {model.Guest, model.Host},
		authServicePath + "DeleteUser": {model.Guest, model.Host},
	}
}
