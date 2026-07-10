package api

import (
	_ "embed"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/rest"

	"zeebee-backend/app/borrower/api/internal/handler"
	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/config"
	"zeebee-backend/common/response"
)

//go:embed borrower.api
var borrowerApiContent string

type ApiServer struct {
	s         *rest.Server
	startOpt  rest.StartOption
}

func NewApiServer(c config.Config) (s service.Service) {
	apiConfig := c.ApiConf
	runOpts := []rest.RunOption{
		rest.WithNotFoundHandler(response.NotFoundHandler{}),
		rest.WithNotAllowedHandler(response.NotAllowedHandler{}),
		rest.WithCors("*"),
	}
	restServer := rest.MustNewServer(apiConfig.RestConf, runOpts...)
	// restServer.Use(middleware.CorsMiddleware)

	// Wrap the http.Server handler ABOVE go-zero's TraceHandler/LogHandler so the
	// built-in access log line carries service_name.
	startOpt := func(svr *http.Server) {
		base := svr.Handler
		svr.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := logx.ContextWithFields(r.Context(),
				logx.Field("service_name", apiConfig.Name))
			base.ServeHTTP(w, r.WithContext(ctx))
		})
	}
	logger := logx.WithCallerSkip(0).
		WithFields(logx.Field("service_name", apiConfig.Name))

	ctx := svc.NewServiceContext(apiConfig, logger, borrowerApiContent)

	// register handlers
	handler.RegisterHandlers(restServer, ctx)

	// extra route
	handler.RegisterExtraRoutes(restServer)

	return &ApiServer{
		s:        restServer,
		startOpt: startOpt,
	}
}

func (s *ApiServer) Start() {
	go func() {
		s.s.StartWithOpts(s.startOpt)
	}()
}

func (s *ApiServer) Stop() {
	s.s.Stop()
}
