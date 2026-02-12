package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"

	"zeebee-backend/app/pdf_generator/api"
	"zeebee-backend/app/pdf_generator/config"
)

var configFile = flag.String("f", "etc/config.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	logx.Infof("load service config successfully")

	apiServer, err := api.NewApiServer(c)
	if err != nil {
		logx.Errorf("start http service failed: %v", err)
		return
	}
	logx.Infof("start http service in %s mode", c.ApiConf.Mode)

	serviceGrp := service.NewServiceGroup()
	serviceGrp.Add(apiServer)
	serviceGrp.Start()
	logx.Infof("api service start success on %d", c.ApiConf.Port)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	fmt.Printf("shutdown service ....")
	serviceGrp.Stop()
}
