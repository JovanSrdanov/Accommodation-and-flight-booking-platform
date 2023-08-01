package interceptor

import (
	"authorization_service/domain/model"
	"context"
	"crypto/x509"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
	"log"
	"regexp"
	"strings"
)

type AuthServerInterceptor struct {
}

func NewAuthServerInterceptor() *AuthServerInterceptor {
	return &AuthServerInterceptor{}
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

	allowedRoles, ok := getProtectedMethodsWithAllowedRoles()[method]
	if !ok {
		// if a provided method is not in the accessible roles map, it means that everyone can use it
		log.Println("Method: " + method + " not found in the list of allowed methods")
		return nil, nil
	}

	clientCert, err := getCertificateFromContext(ctx)
	if err != nil {
		return status.Error(codes.Internal, "failed to get certificate from context"), nil
	}

	providedRole := model.ROLE_UNKNOWN

	// Extract the custom role from the certificate (using the custom OID "1.2.3.4.5").
	for _, ext := range clientCert.Extensions {
		if ext.Id.String() == "1.2.3.4.5" {
			providedRole = getRoleFromExtension(ext.Value)
			//log.Printf("Value of rol: %q", providedRole)
			break
		}
	}

	if providedRole == model.ROLE_UNKNOWN {
		log.Println("NEPOZNAT SERVIS")
		return status.Error(codes.PermissionDenied, "no permission to access this RPC"), nil
	}

	for _, role := range allowedRoles {
		if role == providedRole {
			log.Println("AUTHORIZATION SUCCESSFUL")
			return nil, ctx
		}
	}

	log.Println("NEVALIDNA ROLA")
	return status.Error(codes.PermissionDenied, "no permission to access this RPC"), nil
}

func getRoleFromExtension(extensionValue []byte) model.ServiceRole {
	// Use regular expression to remove non-printable characters
	re := regexp.MustCompile(`[[:cntrl:]]`)
	role := re.ReplaceAllString(strings.TrimSpace(string(extensionValue)), "")
	return convertToEnum(role)
}

func convertToEnum(role string) model.ServiceRole {
	switch role {
	case "ROLE_ACCOMMODATION_SERVICE":
		return model.ROLE_ACCOMMODATION_SERVICE
	case "ROLE_API_GATEWAY":
		return model.ROLE_API_GATEWAY
	case "ROLE_AUTHORIZATION_SERVICE":
		return model.ROLE_AUTHORIZATION_SERVICE
	case "ROLE_NOTIFICATION_SERVICE":
		return model.ROLE_NOTIFICATION_SERVICE
	case "ROLE_RATING_SERVICE":
		return model.ROLE_RATING_SERVICE
	case "ROLE_RESERVATION_SERVICE":
		return model.ROLE_RESERVATION_SERVICE
	case "ROLE_USER_PROFILE_SERVICE":
		return model.ROLE_USER_PROFILE_SERVICE
	default:
		return model.ROLE_UNKNOWN
	}
}

func getCertificateFromContext(ctx context.Context) (*x509.Certificate, error) {
	// Get the peer certificate from the context.
	peerTest, ok := peer.FromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Internal, "Failed to extract peer from context")
	}

	// Extract the client certificate from the peer.
	info, ok := peerTest.AuthInfo.(credentials.TLSInfo)
	if !ok {
		return nil, status.Errorf(codes.Internal, "Failed to extract TLSInfo from peer")
	}

	return info.State.PeerCertificates[0], nil
}

// returns a map which consists of a list of grpc methods and allowed roles for each of them
func getProtectedMethodsWithAllowedRoles() map[string][]model.ServiceRole {
	const authServicePath = "/authorization.AuthorizationService/"

	return map[string][]model.ServiceRole{
		authServicePath + "GetByUsername":  {model.ROLE_API_GATEWAY},
		authServicePath + "ChangeUsername": {model.ROLE_API_GATEWAY},
		authServicePath + "ChangePassword": {model.ROLE_API_GATEWAY},
		authServicePath + "CheckIfDeleted": {model.ROLE_API_GATEWAY},
		authServicePath + "GetById":        {model.ROLE_API_GATEWAY, model.ROLE_USER_PROFILE_SERVICE},
	}
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
