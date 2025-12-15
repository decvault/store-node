package grpcsrv

import (
	"github.com/decvault/library/common/grpcsrv"
	"github.com/decvault/store-node/internal/grpcsrv/options"
	"github.com/decvault/store-node/internal/grpcsrv/store_node"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	api "github.com/decvault/store-node/internal/pb/github.com/decvault/store_node/api"
)

func ProvideGrpcServerSetupFunc() grpcsrv.SetupFunc {
	return func(server *grpc.Server) {
		api.RegisterStoreNodeServer(
			server,
			store_node.NewService(),
		)

		reflection.Register(server)
	}
}

func ProvideGrpcServerSetupOpts(
	options options.Params,
) grpcsrv.SetupOpts {
	return []grpc.ServerOption{
		grpc.UnaryInterceptor(options.PanicInterceptor),
	}
}
