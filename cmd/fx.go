package main

import (
	"context"

	"github.com/decvault/library/badger"
	"github.com/decvault/library/common/config"
	"github.com/decvault/library/common/grpcsrv"
	"github.com/decvault/library/common/logging"
	grpcsetup "github.com/decvault/store-node/internal/grpcsrv"
	"github.com/decvault/store-node/internal/grpcsrv/store_node"
	"github.com/decvault/store-node/internal/pkg/storage"
	"github.com/sirupsen/logrus"
	"go.uber.org/fx"
)

func App() *fx.App {
	return fx.New(
		fx.Options(
			config.Module(),
			logging.Module(),
			grpcsrv.Module(),
			grpcsetup.Module(),
			badger.Module(),
		),
		fx.Provide(
			grpcsetup.NewGrpcServerSetupFunc,
			grpcsetup.NewGrpcServerSetupOpts,
			store_node.NewService,
			storage.NewShardStorage,
		),
		fx.Invoke(func(lc fx.Lifecycle, srv grpcsrv.GrpcServer, logger *logrus.Logger) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					go func() {
						err := srv.Run(ctx)
						if err != nil {
							logger.
								WithContext(ctx).
								Fatalf("failed to start grpc server: %+v", err)
						}
					}()

					return nil
				},
				OnStop: func(ctx context.Context) error {
					srv.GracefulStop(ctx)
					return nil
				},
			})
		}),
	)
}
