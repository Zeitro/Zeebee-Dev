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

	"zeebee-backend/app/borrower/api"
	"zeebee-backend/app/borrower/config"
	"zeebee-backend/app/borrower/rpc"
)

var configFile = flag.String("f", "etc/config.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	logx.Infof("load service config successfully")

	apiServer := api.NewApiServer(c)
	rpcServer, err := rpc.NewRpcServer(c)
	if err != nil {
		fmt.Printf("start rpc server err: %v", err)
		panic(err)
	}
	logx.Infof("start rpc service in %s mode", c.RpcConf.Mode)
	logx.Infof("start http service in %s mode", c.ApiConf.Mode)

	serviceGrp := service.NewServiceGroup()
	serviceGrp.Add(apiServer)
	serviceGrp.Add(rpcServer)
	serviceGrp.Start()
	logx.Infof("api service start success on %d", c.ApiConf.Port)
	logx.Infof("rpc service start success on %s", c.RpcConf.ListenOn)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	fmt.Printf("shutdown service ....")
	serviceGrp.Stop()
	// time.Sleep(5 * time.Second)
}
