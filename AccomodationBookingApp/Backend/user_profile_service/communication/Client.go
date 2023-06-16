package communication

import (
	accommodation "common/proto/accommodation_service/generated"
	authorization "common/proto/authorization_service/generated"
	reservation "common/proto/reservation_service/generated"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"user_profile_service/communication/middleware"
)

func NewAccountCredentialsClient(address string) authorization.AuthorizationServiceClient {
	conn, err := getConnection(address)
	if err != nil {
		log.Fatalf("Failed to start gRPC connection to AccountCredentials service: %v", err)
	}

	return authorization.NewAuthorizationServiceClient(conn)
}

func NewReservationClient(address string) reservation.ReservationServiceClient {
	conn, err := getConnection(address)
	if err != nil {
		log.Fatalf("Failed to start gRPC connection to AccountCredentials service: %v", err)
	}

	return reservation.NewReservationServiceClient(conn)
}

func NewAccommodationClient(address string) accommodation.AccommodationServiceClient {
	conn, err := getConnection(address)
	if err != nil {
		log.Fatalf("Failed to start gRPC connection to AccountCredentials service: %v", err)
	}

	return accommodation.NewAccommodationServiceClient(conn)
}
func getConnection(address string) (*grpc.ClientConn, error) {
	return grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(middleware.NewGRPUnaryClientInterceptor()))
}
