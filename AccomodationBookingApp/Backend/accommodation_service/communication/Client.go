package communication

import (
	"accommodation_service/communication/middleware"
	reservation "common/proto/reservation_service/generated"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

func NewReservationClient(address string) reservation.ReservationServiceClient {
	conn, err := getConnection(address)
	if err != nil {
		log.Fatalf("Failed to start gRPC connection to AccountCredentials service: %v", err)
	}

	return reservation.NewReservationServiceClient(conn)
}

func getConnection(address string) (*grpc.ClientConn, error) {
	return grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(middleware.NewGRPUnaryClientInterceptor()))
}
