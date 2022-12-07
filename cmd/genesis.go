package cmd

import (
	"chain/chain"
	"chain/internal/log"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

// NewGenesisCmd 生成创世块
func NewGenesisCmd(f func(fx.Option)) *cobra.Command {
	return &cobra.Command{
		Use:   "genesis",
		Short: "生成创世块",
		Long:  `生成创世块`,
		Run: func(cmd *cobra.Command, args []string) {
			f(fx.Invoke(func(bc *chain.Blockchain, log log.Logger) {
				fmt.Println("genesis->", args)
				if len(args) == 2 {
					val, _ := strconv.Atoi(args[1])
					bc.CreateGenesisTransaction(args[0], val)
				}
			}))
		},
	}
}
