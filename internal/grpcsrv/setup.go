package grpcsrv

import (
	"github.com/decvault/library/common/grpcsrv"
	panicintc "github.com/decvault/library/common/grpcsrv/options/interceptors/unary/panic"
	"github.com/decvault/store-node/internal/grpcsrv/meta_store"
	"github.com/decvault/store-node/internal/grpcsrv/shard_store"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	meta_store_api "github.com/decvault/store-node/internal/pb/github.com/decvault/meta_store/api"
	shard_store_api "github.com/decvault/store-node/internal/pb/github.com/decvault/shard_store/api"
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

func NewGrpcServerSetupFunc(
	shardStore *shard_store.Service,
	metaStore *meta_store.Service,
) grpcsrv.SetupFunc {
	return func(server *grpc.Server) {
		shard_store_api.RegisterShardStoreServer(
			server,
			shardStore,
		)

		meta_store_api.RegisterMetaStoreServer(
			server,
			metaStore,
		)

		reflection.Register(server)
	}
}
