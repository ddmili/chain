package cmd

import (
	"chain/config"
	"chain/internal/log"

	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

// NewVersion create version cmd
// func NewVersion(cg *config.Config) *cobra.Command {
// 	return &cobra.Command{
// 		Use:   "version",
// 		Short: "打印版本号",
// 		Long:  `打印版本号`,
// 		Run: func(cmd *cobra.Command, args []string) {
// 			fmt.Printf("Git Commit Hash: %s \n", cg.GitHash)
// 			fmt.Printf("Build TimeStamp: %s \n", cg.BuildTime)
// 			fmt.Printf("GoLang Version: %s \n", cg.GoVersion)
// 		},
// 	}
// }

// NewVersion create version cmd
func NewVersion(f func(opt fx.Option)) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "打印版本号",
		Long:  `打印版本号`,
		Run: func(cmd *cobra.Command, args []string) {
			f(fx.Invoke(func(cg *config.Config, lg log.Logger) {
				lg.Info("Git Commit Hash: %s", cg.GitHash)
				lg.Info("Build TimeStamp: %s", cg.BuildTime)
				lg.Info("GoLang Version: %s", cg.GoVersion)
			}))
		},
	}
}
