package handler

import (
	"net/http"
	"net/http/pprof"

	"github.com/zeromicro/go-zero/rest"

	"zeebee-backend/app/admin/api/swagger"
	"zeebee-backend/common/response"
)

func RegisterExtraRoutes(server *rest.Server) {
	registerPprofRoutes(server)
	registerSwaggerRoutes(server)
}

func registerSwaggerRoutes(server *rest.Server) {
	server.AddRoute(rest.Route{
		Method: http.MethodGet,
		Path:   "/swagger",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			// read swagger.json file
			data, err := swagger.SwaggerJSON.ReadFile("admin.json")
			if err != nil {
				response.NotFound(w, "swagger.json not found")
				return
			}

			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write(data)
		},
	})
}

func registerPprofRoutes(server *rest.Server) {

	// pprof routes
	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/debug/pprof/",
		Handler: pprof.Index,
	})
	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/debug/pprof/cmdline",
		Handler: pprof.Cmdline,
	})
	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/debug/pprof/profile",
		Handler: pprof.Profile,
	})
	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/debug/pprof/symbol",
		Handler: pprof.Symbol,
	})
	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/debug/pprof/trace",
		Handler: pprof.Trace,
	})
	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/debug/pprof/allocs",
		Handler: pprof.Handler("allocs").ServeHTTP,
	})
	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/debug/pprof/block",
		Handler: pprof.Handler("block").ServeHTTP,
	})
	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/debug/pprof/goroutine",
		Handler: pprof.Handler("goroutine").ServeHTTP,
	})
	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/debug/pprof/heap",
		Handler: pprof.Handler("heap").ServeHTTP,
	})
	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/debug/pprof/mutex",
		Handler: pprof.Handler("mutex").ServeHTTP,
	})
	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/debug/pprof/threadcreate",
		Handler: pprof.Handler("threadcreate").ServeHTTP,
	})

}
