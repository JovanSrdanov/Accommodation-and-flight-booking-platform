package handler

import "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

type IHandler interface {
	Init(mux *runtime.ServeMux)
}
