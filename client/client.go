package client

import (
	"chain/chain"
	"chain/cmd"
	"chain/config"
	"chain/internal/log"
	"chain/internal/store"
	"chain/internal/wallet"
	"chain/p2p"
	"context"
	"os"
	"os/signal"

	"go.uber.org/fx"
)

var (
	gitHash   string
	buildTime string
	goVersion string
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	kill := make(chan os.Signal, 1)
	signal.Notify(kill)

	go func() {
		<-kill
		cancel()
	}()

	app := fx.New(
		fx.Provide(config.NewCg),
		fx.Provide(log.NewLogger),
		fx.Provide(store.NewDb),
		fx.Provide(p2p.NewNetwork),
		fx.Provide(wallet.NewWallets),
		fx.Provide(chain.NewBlockchain),
		fx.Invoke(setVersion),
		fx.Invoke(cmd.Execute),
	)
	err := app.Start(ctx)
	if err != nil {
		panic(err)
	}

}

// setVersion 设置版本号
func setVersion(cg *config.Config) {
	cg.GitHash = gitHash
	cg.BuildTime = gitHash
	cg.GoVersion = goVersion
	// fmt.Printf("config:%+v", cg)
}
