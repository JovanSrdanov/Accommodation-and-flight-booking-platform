package startup

import (
	"authorization_service/domain/model"
	"authorization_service/domain/token"
	"authorization_service/interceptor"
	notification "common/proto/notification_service/generated"

	"fmt"
	"google.golang.org/grpc"
	"gorm.io/gorm"
	"log"
	"net"
	"notification_service/communication/handler"
	"notification_service/domain/service"
	"notification_service/persistence/repository"
)

type Server struct {
	config *Configuration
}

func NewServer(config *Configuration) *Server {
	return &Server{config: config}
}

func (server *Server) Start() {
	postgresClient := server.initPostgresClient()
	notificationConsentRepo := initNotificationConsentRepo(postgresClient)
	notificationConsentService := service.NewNotificationConsentService(*notificationConsentRepo)
	notificationConsentHandler := handler.NewNotificationConsentHandler(*notificationConsentService)

	server.startGrpcServer(notificationConsentHandler)
}

func initNotificationConsentRepo(client *gorm.DB) *repository.NotificationConsentRepositoryPG {
	repo, err := repository.NewNotificationConsentServicePG(client)
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
func (server *Server) startGrpcServer(notificationConsentHandler *handler.NotificationConsentHandler) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", server.config.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	var tokenMaker, _ = token.NewPasetoMaker("12345678901234567890123456789012")
	protectedMethodsWithAllowedRoles := getProtectedMethodsWithAllowedRoles()
	authInterceptor := interceptor.NewAuthServerInterceptor(tokenMaker, protectedMethodsWithAllowedRoles)

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(authInterceptor.Unary()),
		grpc.StreamInterceptor(authInterceptor.Stream()),
	)
	notification.RegisterNotificationServiceServer(grpcServer, notificationConsentHandler)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}

// returns a map which consists of a list of grpc methods and allowed roles for each of them
func getProtectedMethodsWithAllowedRoles() map[string][]model.Role {
	const authServicePath = "/notification.NotificationService/"

	return map[string][]model.Role{
		authServicePath + "UpdateMyNotificationConsent": {model.Guest, model.Host},
		authServicePath + "GetMyNotificationSettings":   {model.Guest, model.Host},
	}
}
