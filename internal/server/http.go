package server

import (
	v1 "github.com/YangZhaoWeblog/UserService/api/helloworld/v1"
	userv1 "github.com/YangZhaoWeblog/UserService/api/user/v1"
	"github.com/YangZhaoWeblog/UserService/internal/conf"
	"github.com/YangZhaoWeblog/UserService/internal/pkg"
	"github.com/YangZhaoWeblog/UserService/internal/service"
	"github.com/go-kratos/kratos/v2/middleware/metrics"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server, greeter *service.GreeterService, user *service.UserService, metricsData *pkg.MetricsData) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
			tracing.Server(),
			metrics.Server(
				metrics.WithSeconds(metricsData.Seconds),
				metrics.WithRequests(metricsData.Requests),
			),
		),
	}

	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}

	srv := http.NewServer(opts...)

	srv.Handle("/metrics", promhttp.Handler())

	v1.RegisterGreeterHTTPServer(srv, greeter)
	userv1.RegisterUserHTTPServer(srv, user)
	return srv
}
