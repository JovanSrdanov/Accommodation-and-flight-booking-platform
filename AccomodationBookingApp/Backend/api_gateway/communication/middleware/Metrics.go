package middleware

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

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

var (
	UniqueUserCount = promauto.NewCounter(prometheus.CounterOpts{
		Name: "unique_user_visits",
		Help: "The total number visits from unique users",
	})
)
