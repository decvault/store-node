package grpcsrv

import (
	"github.com/decvault/library/common/grpcsrv"
	panicintc "github.com/decvault/library/common/grpcsrv/options/interceptors/unary/panic"
	"github.com/decvault/store-node/internal/grpcsrv/store_node"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	api "github.com/decvault/store-node/internal/pb/github.com/decvault/store_node/api"
)

type Params struct {
	fx.In

	PanicInterceptor grpc.UnaryServerInterceptor `name:"panic_handler"`
}

func Module() fx.Option {
	return fx.Module(
		"grpcsrv_opts",
		panicintc.Module(),
	)
}

func NewGrpcServerSetupOpts(opts Params) grpcsrv.SetupOpts {
	return []grpc.ServerOption{
		grpc.UnaryInterceptor(opts.PanicInterceptor),
	}
}

func NewGrpcServerSetupFunc(service *store_node.Service) grpcsrv.SetupFunc {
	return func(server *grpc.Server) {
		api.RegisterStoreNodeServer(
			server,
			service,
		)

		reflection.Register(server)
	}
}
