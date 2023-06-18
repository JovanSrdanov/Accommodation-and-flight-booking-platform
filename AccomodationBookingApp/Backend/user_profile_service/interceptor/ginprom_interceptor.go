package interceptor

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"log"
)

type GinpromInterceptor struct{}

func NewGinpromInterceptor() *GinpromInterceptor {
	return &GinpromInterceptor{}
}

func (interceptor *GinpromInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		unaryHandler grpc.UnaryHandler,
	) (interface{}, error) {
		log.Println("--> unary ginprom interceptor: ", info.FullMethod)

		err, newCtx := interceptor.instrument(ctx)
		if err != nil {
			return nil, err
		}

		return unaryHandler(newCtx, req)
	}
}

func (interceptor *GinpromInterceptor) instrument(ctx context.Context) (error, context.Context) {
	log.Println("Instrumentation in progress...")

	metaData, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		log.Println("metadata is not provided")
		//return status.Errorf(codes.Aborted, "metadata is not provided"), nil
	}

	log.Println("Metadata: ")
	log.Println(metaData)

	values := metaData["ginprom"]
	if len(values) == 0 {
		log.Println("ginprom instance not provided")
		//return status.Errorf(codes.Aborted, "ginprom instance not provided"), nil
	}

	log.Println("Ginprom instance: ")
	log.Println(values)

	ctx = context.WithValue(ctx, "ginprom", values)

	return nil, ctx
}
