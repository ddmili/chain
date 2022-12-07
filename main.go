package main

import (
	"chain/chain"
	"chain/client"
	"chain/cmd"
	"chain/config"
	"chain/internal/log"
	"chain/internal/store"
	"chain/internal/wallet"
	"chain/p2p"
	"chain/service"
	"context"

	"go.uber.org/fx"
)

var (
	gitHash   string
	buildTime string
	goVersion string
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	opts := []fx.Option{}
	opts = append(opts, fx.Provide(ctx, cancel))
	opts = append(opts,
		fx.Provide(config.NewCg),
		fx.Provide(log.NewLogger),
		fx.Provide(store.NewDb),
		fx.Provide(p2p.NewNetwork),
		fx.Provide(wallet.NewWallets),
		fx.Provide(chain.NewBlockchain),
		fx.Provide(service.NewChainService),
		fx.Provide(client.NewChainClient),
		fx.Invoke(setVersion),
		// fx.NopLogger,
	)
	opts = append(opts, cmd.Execute()...)

	fx.New(opts...)

}

// setVersion 设置版本号
func setVersion(cg *config.Config) {
	cg.GitHash = gitHash
	cg.BuildTime = gitHash
	cg.GoVersion = goVersion
	// fmt.Printf("config:%+v", cg)
}
