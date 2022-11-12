package cmd

import (
	"chain/config"
	"chain/internal/log"
	"chain/service"

	"github.com/spf13/cobra"
	"go.uber.org/fx"
	"golang.org/x/sync/errgroup"
)

// NewStartCmd 开启服务
func NewStartCmd(f func(fx.Option)) *cobra.Command {
	return &cobra.Command{
		Use:   "start",
		Short: "start chain node",
		Long:  `start chain node`,
		Run: func(cmd *cobra.Command, args []string) {
			f(fx.Invoke(func(lf fx.Lifecycle, rpcHTTP *service.ChainService, lg log.Logger, cg *config.Config) {
				var srv errgroup.Group
				srv.Go(func() error {
					return service.GrpcServer(rpcHTTP, lg, cg)
				})
				if err := srv.Wait(); err != nil {
					log.Fatal("start cmd run error:%s", err)
				}
			}))

		},
	}
}
