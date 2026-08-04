package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/tenant/config"
	"zeebee-backend/app/tenant/rpc"
)

var configFile = flag.String("f", "etc/tenant.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	logx.Infof("load tenant service config successfully")

	rpcServer, err := rpc.NewRpcServer(c)
	if err != nil {
		panic(err)
	}
	logx.Infof("start tenant rpc service in %s mode on %s", c.RpcConf.Mode, c.RpcConf.ListenOn)

	go func() {
		rpcServer.Start()
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	fmt.Printf("shutdown tenant service ....")
	rpcServer.Stop()
}
