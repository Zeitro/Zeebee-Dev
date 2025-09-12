package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/zeromicro/go-zero/core/conf"
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

	apiServer := api.NewApiServer(c)
	rpcServer, err := rpc.NewRpcServer(c)
	if err != nil {
		panic(err)
	}

	serviceGrp := service.NewServiceGroup()
	serviceGrp.Add(apiServer)
	serviceGrp.Add(rpcServer)
	serviceGrp.Start()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	fmt.Printf("shutdown service ....")
	serviceGrp.Stop()
	// time.Sleep(5 * time.Second)
}
