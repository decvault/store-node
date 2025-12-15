package options

import (
	"github.com/decvault/library/common/grpcsrv/options/interceptors/unary/panic"
	"go.uber.org/fx"
	"google.golang.org/grpc"
)

func Module() fx.Option {
	return fx.Module(
		"grpcsrv_opts",
		fx.Options(
			panicintc.Module(),
		),
	)
}

type Params struct {
	fx.In

	PanicInterceptor grpc.UnaryServerInterceptor `name:"panic_handler"`
}
