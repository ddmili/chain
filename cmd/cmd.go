package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

// Execute creates a command
func Execute() []fx.Option {
	opts := []fx.Option{}
	var rootCmd = &cobra.Command{Use: "chain_game"}
	optFun := func(opt fx.Option) {
		opts = append(opts, opt)
	}
	rootCmd.AddCommand(NewVersion(optFun))
	rootCmd.AddCommand(NewGenesisCmd(optFun))
	rootCmd.AddCommand(NewWallets(optFun))
	rootCmd.AddCommand(NewStartCmd(optFun))
	err := rootCmd.Execute()
	// fmt.Println("start root cmd")
	if err != nil {
		fmt.Printf("start root cmd error: %s", err)
		rootCmd.Help()
	}
	return opts
}
