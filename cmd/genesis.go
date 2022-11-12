package cmd

import (
	"chain/chain"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

// NewGenesisCmd 生成创世块
func NewGenesisCmd(bc *chain.Blockchain) *cobra.Command {
	return &cobra.Command{
		Use:   "genesis",
		Short: "生成创世块",
		Long:  `生成创世块`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("genesis->", args)
			if len(args) == 2 {
				val, _ := strconv.Atoi(args[1])
				bc.CreateGenesisTransaction(args[0], val)
			}
		},
	}
}
