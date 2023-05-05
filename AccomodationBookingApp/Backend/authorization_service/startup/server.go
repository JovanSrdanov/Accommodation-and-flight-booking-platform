package startup

import (
	"authorization_service/communication/handler"
	"authorization_service/domain/model"
	"authorization_service/domain/service"
	"authorization_service/domain/token"
	"authorization_service/interceptor"
	"authorization_service/persistence/repository"
	"authorization_service/startup/configuration"
	authorization "common/proto/authorization_service/generated"
	"fmt"
	"google.golang.org/grpc"
	"gorm.io/gorm"
	"log"
	"net"
)

type Server struct {
	config *configuration.Configuration
}

func NewServer(config *configuration.Configuration) *Server {
	return &Server{config: config}
}

func (server Server) Start() {
	postgresClient := server.initPostgresClient()
	accountCredentialsRepo := server.initAccountCredentialsRepo(postgresClient)

	// TODO Stefan: currently not working with .env file
	tokenMaker, err := token.NewPasetoMaker("12345678901234567890123456789012")
	if err != nil {
		log.Fatalf("cannot create token maker: %v", err)
	}

	accountCredentialsService := service.NewAccountCredentialsService(accountCredentialsRepo, tokenMaker)
	accountCredentialsHandler := handler.NewAccountCredentialsHandler(accountCredentialsService)

	server.startGrpcServer(accountCredentialsHandler, tokenMaker)
}

func (server Server) initPostgresClient() *gorm.DB {
	client, err := repository.GetClient(
		server.config.DBHost, server.config.DBUser,
		server.config.DBPass, server.config.DBName,
		server.config.DBPort)
	if err != nil {
		log.Fatal(err)
	}
	return client
}

func (server Server) initAccountCredentialsRepo(postgresClient *gorm.DB) *repository.AccountCredentialsRepositoryPG {
	repo, err := repository.NewAccountCredentialsRepositoryPG(postgresClient)
	if err != nil {
		log.Fatal(err)
	}
	return repo
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
	// TODO Stefan add real method map
	tempRoles := make(map[string][]model.Role)
	authInterceptor := interceptor.NewAuthServerInterceptor(maker, tempRoles)

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(authInterceptor.Unary()),
		grpc.StreamInterceptor(authInterceptor.Stream()),
	)
	authorization.RegisterAuthorizationServiceServer(grpcServer, accountCredentialsHandler)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
