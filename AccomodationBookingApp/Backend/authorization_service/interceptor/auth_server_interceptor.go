package interceptor

import (
	"authorization_service/domain/model"
	"authorization_service/domain/token"
	"context"
	"log"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type AuthServerInterceptor struct {
	// TODO Stefan check if should be *token.PasetoMaker
	tokenMaker                       token.Maker
	protectedMethodsWithAllowedRoles map[string][]model.Role
}

func NewAuthServerInterceptor(tokenMaker token.Maker, accessibleRoles map[string][]model.Role) *AuthServerInterceptor {
	return &AuthServerInterceptor{tokenMaker: tokenMaker, protectedMethodsWithAllowedRoles: accessibleRoles}
}

func (interceptor *AuthServerInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		unaryHandler grpc.UnaryHandler,
	) (interface{}, error) {
		log.Println("--> unary interceptor: ", info.FullMethod)

		err, newCtx := interceptor.authorize(ctx, info.FullMethod)
		if err != nil {
			return nil, err
		}

		return unaryHandler(newCtx, req)
	}
}

func (interceptor *AuthServerInterceptor) Stream() grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		stream grpc.ServerStream,
		info *grpc.StreamServerInfo,
		streamHandler grpc.StreamHandler,
	) error {
		log.Println("--> stream interceptor: ", info.FullMethod)

		err, newCtx := interceptor.authorize(stream.Context(), info.FullMethod)
		if err != nil {
			return err
		}

		wrappedStream := newWrappedServerStream(stream, newCtx)

		return streamHandler(srv, wrappedStream)
	}
}

func (interceptor *AuthServerInterceptor) authorize(ctx context.Context, method string) (error, context.Context) {

	allowedRoles, ok := interceptor.protectedMethodsWithAllowedRoles[method]
	if !ok {
		// if a provided method is not in the accessible roles map, it means that everyone can use it
		log.Println("Method: " + method + " not found in the list of allowed methods")
		return nil, nil
	}
	metaData, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return status.Errorf(codes.Unauthenticated, "metadata is not provided"), nil
	}

	values := metaData["authorization"]
	if len(values) == 0 {
		return status.Errorf(codes.Unauthenticated, "authorization token not provided"), nil
	}

	log.Println("Header")

	accessToken := strings.TrimPrefix(values[0], "Bearer ")
	tokenPayload, err := interceptor.tokenMaker.VerifyToken(accessToken)
	if err != nil {
		return status.Errorf(codes.Unauthenticated, "access token is invalid: ", err), nil
	}

	ctx = context.WithValue(ctx, "id", tokenPayload.ID)

	providedRole := tokenPayload.Role

	for _, role := range allowedRoles {
		if role == providedRole {
			return nil, ctx
		}
	}

	return status.Error(codes.PermissionDenied, "no permission to access this RPC"), nil
}

type wrappedServerStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (w *wrappedServerStream) Context() context.Context {
	return w.ctx
}

func newWrappedServerStream(stream grpc.ServerStream, ctx context.Context) grpc.ServerStream {
	return &wrappedServerStream{
		ServerStream: stream,
		ctx:          ctx,
	}
}
