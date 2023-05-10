package startup

import (
	"authorization_service/domain/model"
	"authorization_service/domain/token"
	"authorization_service/interceptor"
	reservation "common/proto/reservation_service/generated"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
	"log"
	"net"
	"reservation_service/communication/handler"
	"reservation_service/domain/service"
	"reservation_service/persistence/repository"
)

type Server struct {
	config *Configuration
}

func NewServer(config *Configuration) *Server {
	return &Server{config: config}
}

func (server Server) Start() {
	mongoClient := server.initMongoClient()
	reservationRepo := initUserProfileRepo(mongoClient)
	reservationService := service.NewReservationService(*reservationRepo)
	reservationHandler := handler.NewReservationHandler(*reservationService)
	server.startGrpcServer(reservationHandler)
}

func (server Server) initMongoClient() *mongo.Client {
	client, err := repository.GetClient(server.config.DBPort, server.config.DBName)
	if err != nil {
		log.Fatal(err)
	}
	return client
}

func initUserProfileRepo(mongoClient *mongo.Client) *repository.ReservationRepositoryMongo {
	repo, err := repository.NewReservationRepositoryMongo(mongoClient)
	if err != nil {
		log.Fatal(err)
	}
	return repo
}
func (server Server) startGrpcServer(reservationHandler *handler.ReservationHandler) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", server.config.Port))
	log.Println("port: " + server.config.Port)
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
	reservation.RegisterReservationServiceServer(grpcServer, reservationHandler)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}

// returns a map which consists of a list of grpc methods and allowed roles for each of them
func getProtectedMethodsWithAllowedRoles() map[string][]model.Role {
	const authServicePath = "/reservation.ReservationService/"

	return map[string][]model.Role{
		authServicePath + "AcceptReservation": {model.Guest},
	}
}
