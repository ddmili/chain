package p2p

import (
	"chain/config"
	"context"
	"fmt"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/multiformats/go-multiaddr"
	"go.uber.org/fx"
)

// Network p2p network
type Network struct {
	host host.Host
}

// NewNetwork create a p2p network
func NewNetwork(cg *config.Config, lc fx.Lifecycle) *Network {
	if cg.ListenAddresses == "" {
		return &Network{}
	}
	listenAddress, err := multiaddr.NewMultiaddr(cg.ListenAddresses)
	if err != nil {
		panic(err)
	}
	node, err := libp2p.New(libp2p.ListenAddrs(listenAddress))
	if err != nil {
		panic(err)
	}
	// print the node's listening addresses
	fmt.Println("Listen addresses:", node.Addrs())

	net := &Network{host: node}
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			fmt.Println("shutdown p2p server")
			err := net.Close()
			if err != nil {
				fmt.Println("p2p shutdown error:", err)
			}
			return err
		},
	})
	return net
}

// Close close network
func (n *Network) Close() error {
	err := n.host.Close()
	if err != nil {
		fmt.Println("p2p shutdown error:", err)
	}
	return err
}
