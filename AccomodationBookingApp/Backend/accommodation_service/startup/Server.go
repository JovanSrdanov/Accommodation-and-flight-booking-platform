package startup

import (
	"accommodation_service/communication/handler"
	"accommodation_service/domain/service"
	"accommodation_service/persistence/repository"
	"authorization_service/domain/model"
	"authorization_service/domain/token"
	"authorization_service/interceptor"
	accommodation "common/proto/accommodation_service/generated"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
	"log"
	"net"
)

type Server struct {
	config *Configuration
}

func NewServer(config *Configuration) *Server {
	return &Server{config: config}
}

func (server Server) Start() {
	mongoClient := server.initMongoClient()
	accommodationRepo := initUserProfileRepo(mongoClient)
	accommodationService := service.NewAccommodationService(*accommodationRepo)
	accommodationHandler := handler.NewAccommodationHandler(*accommodationService)
	server.startGrpcServer(accommodationHandler)
}

func (server Server) initMongoClient() *mongo.Client {
	client, err := repository.GetClient(server.config.DbUser, server.config.DbPass)
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
func (server Server) startGrpcServer(userProfileHandler *handler.AccommodationHandler) {
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
	accommodation.RegisterAccommodationServiceServer(grpcServer, userProfileHandler)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}

// returns a map which consists of a list of grpc methods and allowed roles for each of them
func getProtectedMethodsWithAllowedRoles() map[string][]model.Role {
	const authServicePath = "/accommodation.AccommodationService/"

	return map[string][]model.Role{
		authServicePath + "Update": {model.Guest},
	}
}
