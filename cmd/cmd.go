package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

// Execute creates a command
// func Execute(lf fx.Lifecycle, rpcHTTP *service.ChainService, cg *config.Config, bc *chain.Blockchain, wallets *wallet.Wallets, lg log.Logger) {
// 	var rootCmd = &cobra.Command{Use: "chain_game"}
// 	rootCmd.AddCommand(NewVersion(cg))
// 	rootCmd.AddCommand(NewGenesisCmd(bc))
// 	rootCmd.AddCommand(NewWallets(wallets))
// 	rootCmd.AddCommand(NewStartCmd(lf, rpcHTTP, lg, cg))
// 	err := rootCmd.Execute()
// 	if err != nil {
// 		lg.Error("start root cmd error: %s", err)
// 		rootCmd.Help()
// 	}

// 	lg.Info("start root cmd")

// }

// Execute creates a command
func Execute() []fx.Option {
	opts := []fx.Option{}
	var rootCmd = &cobra.Command{Use: "chain_game"}
	rootCmd.AddCommand(NewVersion(func(opt fx.Option) {
		opts = append(opts, opt)
	}))
	rootCmd.AddCommand(NewStartCmd(func(opt fx.Option) {
		opts = append(opts, opt)
	}))
	err := rootCmd.Execute()
	fmt.Printf("start root cmd")
	if err != nil {
		fmt.Printf("start root cmd error: %s", err)
		rootCmd.Help()
	}
	return opts
}
