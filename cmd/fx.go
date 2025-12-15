package main

import (
	"context"

	"github.com/decvault/library/common/config"
	"github.com/decvault/library/common/grpcsrv"
	"github.com/decvault/library/common/logging"
	grpcsetup "github.com/decvault/store-node/internal/grpcsrv"
	grpcopts "github.com/decvault/store-node/internal/grpcsrv/options"
	"github.com/sirupsen/logrus"
	"go.uber.org/fx"
)

func App() *fx.App {
	return fx.New(
		fx.Options(
			config.Module(),
			logging.Module(),
			grpcsrv.Module(),
			grpcopts.Module(),
		),
		fx.Provide(
			grpcsetup.ProvideGrpcServerSetupFunc,
			grpcsetup.ProvideGrpcServerSetupOpts,
		),
		fx.Invoke(func(lc fx.Lifecycle, srv grpcsrv.GrpcServer, logger *logrus.Logger) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					go func() {
						err := srv.Run()
						if err != nil {
							logger.
								WithContext(ctx).
								Fatalf("failed to start grpc server: %+v", err)
						}
					}()

					return nil
				},
				OnStop: func(ctx context.Context) error {
					logger.
						WithContext(ctx).
						Info("shutting down app...")

					return nil
				},
			})
		}),
	)
}
