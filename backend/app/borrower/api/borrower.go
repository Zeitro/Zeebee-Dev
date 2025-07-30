package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"zeebee-backend/app/borrower/api/internal/config"
	"zeebee-backend/app/borrower/api/internal/handler"
	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/common/response"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/borrower-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	runOpts := []rest.RunOption{
		rest.WithNotFoundHandler(response.NotFoundHandler{}),
		rest.WithNotAllowedHandler(response.NotFoundHandler{}),
	}

	server := rest.MustNewServer(c.RestConf, runOpts...)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	server.AddRoute(rest.Route{
		Method: http.MethodGet,
		Path:   "/swagger",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			// 读取 swagger.json 文件
			data, err := os.ReadFile("swagger/borrower.json")
			if err != nil {
				http.Error(w, "swagger.json not found", http.StatusNotFound)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.Write(data)
		},
	})

	go func() {
		fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
		server.Start()
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	fmt.Printf("shutdown service in 5 seconds ....")
	time.Sleep(5 * time.Second)

}
