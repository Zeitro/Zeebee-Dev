package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"zeebee-backend/app/borrower/rpc/internal/config"
	"zeebee-backend/app/borrower/rpc/internal/global"
	applicantServer "zeebee-backend/app/borrower/rpc/internal/server/applicant"
	commonServer "zeebee-backend/app/borrower/rpc/internal/server/common"
	loanServer "zeebee-backend/app/borrower/rpc/internal/server/loan"
	propertyServer "zeebee-backend/app/borrower/rpc/internal/server/property"
	userServer "zeebee-backend/app/borrower/rpc/internal/server/user"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/user.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	// init global config
	err := global.InitGlobal(c)
	if err != nil {
		panic(err)
	}

	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		borrower.RegisterUserServer(grpcServer, userServer.NewUserServer(ctx))
		borrower.RegisterLoanServer(grpcServer, loanServer.NewLoanServer(ctx))
		borrower.RegisterCommonServer(grpcServer, commonServer.NewCommonServer(ctx))
		borrower.RegisterPropertyServer(grpcServer, propertyServer.NewPropertyServer(ctx))
		borrower.RegisterApplicantServer(grpcServer, applicantServer.NewApplicantServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()
	go func() {
		fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
		s.Start()
	}()

	// mercy shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	fmt.Printf("shutdown service in 5 seconds ....")
	time.Sleep(5 * time.Second)
}
