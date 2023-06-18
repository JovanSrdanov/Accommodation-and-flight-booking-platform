package middleware

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

//func RegisterMetrics() (*grpcprom.ClientMetrics, *prometheus.Registry) {
//	// Setup counter metrics.
//	countMetrics := grpcprom.NewClientMetrics(
//		grpcprom.WithClientCounterOptions(grpcprom.WithConstLabels(prometheus.Labels{"test": "test"})),
//	)
//	// grpcprom:Package prometheus provides a standalone interceptor for metrics.
//
//	reg := prometheus.NewRegistry() // NewRegistry creates a new vanilla Registry without any Collectors pre-registered.
//	reg.MustRegister(countMetrics)  // MustRegister implements Registerer
//
//	return countMetrics, reg
//}

var (
	HttpReqCountTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "api_gateway_http_requests_total",
		Help: "The total number of processed http requests",
	})
)

var (
	HttpReqCountSucc = promauto.NewCounter(prometheus.CounterOpts{
		Name: "api_gateway_http_requests_successful",
		Help: "The total number of processed successful http requests",
	})
)

var (
	HttpReqCountFail = promauto.NewCounter(prometheus.CounterOpts{
		Name: "api_gateway_http_requests_failed",
		Help: "The total number of processed failed http requests",
	})
)

var (
	HttpReqCountNotFound = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "api_gateway_http_requests_not_found",
		Help: "The total number of processed failed http requests that produce 404 error",
	},
		[]string{"endpoint"})
)
