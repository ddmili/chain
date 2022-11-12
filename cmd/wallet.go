package cmd

import (
	"chain/internal/wallet"
	"fmt"

	"github.com/spf13/cobra"
)

// NewWallets 钱包操作
func NewWallets(wallets *wallet.Wallets) *cobra.Command {
	c := &cobra.Command{
		Use:   "wallet",
		Short: "钱包操作",
		Long:  `创建钱包,查看钱包,列表钱包`,
	}
	c.AddCommand(generateWallet(wallets))
	return c
}

// generateWallet 创建钱包
func generateWallet(wallets *wallet.Wallets) *cobra.Command {
	return &cobra.Command{
		Use:   "generate",
		Short: "生成钱包",
		Long:  `生成钱包`,
		Run: func(cmd *cobra.Command, args []string) {
			address, privkey, mnemonicWord := wallets.GenerateWallet()
			fmt.Println("助记词：", mnemonicWord)
			fmt.Println("私钥：", privkey)
			fmt.Println("地址：", address)
		},
	}
}
