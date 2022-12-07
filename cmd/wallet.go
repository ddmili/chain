package cmd

import (
	"chain/internal/wallet"
	"fmt"

	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

// NewWallets 钱包操作
func NewWallets(f func(fx.Option)) *cobra.Command {
	c := &cobra.Command{
		Use:   "wallet",
		Short: "钱包操作",
		Long:  `创建钱包,查看钱包,列表钱包`,
	}
	c.AddCommand(generateWallet(f))
	c.AddCommand(walletList(f))
	c.AddCommand(getWallet(f))
	return c
}

// generateWallet 创建钱包
func generateWallet(f func(fx.Option)) *cobra.Command {
	return &cobra.Command{
		Use:   "generate",
		Short: "生成钱包",
		Long:  `生成钱包`,
		Run: func(cmd *cobra.Command, args []string) {
			f(fx.Invoke(func(wallets *wallet.Wallets) {
				address, privkey, mnemonicWord := wallets.GenerateWallet()
				fmt.Println("助记词：", mnemonicWord)
				fmt.Println("私钥：", privkey)
				fmt.Println("地址：", address)
			}))

		},
	}
}

// walletList 钱包列表
func walletList(f func(fx.Option)) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "钱包列表",
		Run: func(cmd *cobra.Command, args []string) {
			f(fx.Invoke(func(wallets *wallet.Wallets) {
				list := wallets.GetAllWallet()
				for _, w := range list {
					fmt.Println("----------------------")
					fmt.Println("助记词：", w.MnemonicWord)
					fmt.Println("私钥：", w.GetPrivateKey())
					fmt.Println("地址：", string(w.GetAddress()))
				}
			}))

		},
	}
}

// getWallet 钱包列表
func getWallet(f func(fx.Option)) *cobra.Command {
	return &cobra.Command{
		Use:   "get",
		Short: "get wallet detail",
		Run: func(cmd *cobra.Command, args []string) {
			f(fx.Invoke(func(wallets *wallet.Wallets) {
				if len(args) == 1 {
					w, err := wallets.GetWallet(args[0])
					if err == nil {
						fmt.Println("助记词：", w.MnemonicWord)
						fmt.Println("私钥：", w.GetPrivateKey())
						fmt.Println("地址：", string(w.GetAddress()))
					}

				}

			}))

		},
	}
}
