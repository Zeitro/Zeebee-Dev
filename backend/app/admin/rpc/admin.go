package rpc

import (
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/zeromicro/go-zero/core/service"

	apiServer "zeebee-backend/app/admin/rpc/internal/server/api"
	organizationUserServer "zeebee-backend/app/admin/rpc/internal/server/organizationuser"
	resourceServer "zeebee-backend/app/admin/rpc/internal/server/resource"
	roleServer "zeebee-backend/app/admin/rpc/internal/server/role"
	userServer "zeebee-backend/app/admin/rpc/internal/server/user"
	"zeebee-backend/app/admin/rpc/pb/admin"

	"zeebee-backend/app/admin/config"
	"zeebee-backend/app/admin/rpc/internal/global"
	"zeebee-backend/app/admin/rpc/internal/svc"
)

type RpcServer struct {
	s *zrpc.RpcServer
}

func NewRpcServer(c config.Config) (s service.Service, err error) {
	rpcConfig := c.RpcConf

	err = global.InitGlobal(rpcConfig)
	if err != nil {
		return nil, err
	}

	ctx := svc.NewServiceContext(rpcConfig)
	rpcServer := zrpc.MustNewServer(rpcConfig.RpcServerConf, func(grpcServer *grpc.Server) {
		admin.RegisterUserServer(grpcServer, userServer.NewUserServer(ctx))
		admin.RegisterApiServer(grpcServer, apiServer.NewApiServer(ctx))
		admin.RegisterResourceServer(grpcServer, resourceServer.NewResourceServer(ctx))
		admin.RegisterRoleServer(grpcServer, roleServer.NewRoleServer(ctx))
		admin.RegisterOrganizationUserServer(grpcServer, organizationUserServer.NewOrganizationUserServer(ctx))
		if rpcConfig.Mode == service.DevMode || rpcConfig.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	return &RpcServer{
		s: rpcServer,
	}, nil
}

func (s *RpcServer) Start() {
	go func() {
		s.s.Start()
	}()
}

func (s *RpcServer) Stop() {
	s.s.Stop()
}
