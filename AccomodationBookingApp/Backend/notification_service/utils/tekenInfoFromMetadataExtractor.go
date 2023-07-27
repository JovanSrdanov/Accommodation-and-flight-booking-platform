package utils

import (
	"authorization_service/domain/token"
	"context"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"strings"
)

func GetTokenInfo(ctx context.Context) (uuid.UUID, error) {
	metaData, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return uuid.Nil, status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}

	values := metaData["authorization"]
	if len(values) == 0 {
		return uuid.Nil, status.Errorf(codes.Unauthenticated, "token info not provided")
	}

	accessToken := strings.TrimPrefix(values[0], "Bearer ")

	tokenMaker, _ := token.NewPasetoMaker("12345678901234567890123456789012")

	tokenPayload, err := tokenMaker.VerifyToken(accessToken)
	if err != nil {
		return uuid.Nil, status.Errorf(codes.Unauthenticated, "access token is invalid: ", err)
	}

	return tokenPayload.ID, nil
}
