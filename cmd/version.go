package cmd

import (
	"chain/client"
	"chain/config"
	"chain/internal/log"
	"chain/service/pd"
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

// NewVersion create version cmd
func NewVersion(f func(opt fx.Option)) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "打印版本号",
		Long:  `打印版本号`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Version: ")
			f(fx.Invoke(func(cli *client.ChainCli, cg *config.Config, lg log.Logger) {
				req := &pd.VersionRequest{}
				v, err := cli.Version(context.Background(), req)
				if err != nil {
					log.Error("get version:%s", err)
					return
				}
				lg.Info(v.Version)
			}))
		},
	}
}
