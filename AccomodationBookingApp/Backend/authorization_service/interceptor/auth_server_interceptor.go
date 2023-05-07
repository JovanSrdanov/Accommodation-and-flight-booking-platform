package interceptor

import (
	"authorization_service/domain/model"
	"authorization_service/domain/token"
	"context"
	"github.com/o1egl/paseto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"log"
	"strings"
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
	log.Println("Authorization in progress...")
	allowedRoles, ok := interceptor.protectedMethodsWithAllowedRoles[method]
	if !ok {
		// if a provided method is not in the accessible roles map, it means that everyone can use it
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

	accessToken := strings.TrimPrefix(values[0], "Bearer ")
	tokenPayload, err := interceptor.tokenMaker.VerifyToken(accessToken)
	if err != nil {
		return status.Errorf(codes.Unauthenticated, "access token is invalid: ", err), nil
	}

	ctx = context.WithValue(ctx, "id", tokenPayload.ID)

	var footerData map[string]interface{}
	if err := paseto.ParseFooter(accessToken, &footerData); err != nil {
		return status.Errorf(codes.Internal, "failed to parse token footer: ", err), nil
	}

	providedRole := int8(footerData["Role"].(float64))
	for _, role := range allowedRoles {
		if int8(role) == providedRole {
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
