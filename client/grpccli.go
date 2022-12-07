package client

import (
	"chain/config"
	"chain/service/pd"
	"log"

	"google.golang.org/grpc"
)

// ChainCli chain client
type ChainCli struct {
	pd.ChainClient
	conn *grpc.ClientConn
}

// NewChainClient creates a new ChainCli
func NewChainClient(cg *config.Config) *ChainCli {
	conn, err := grpc.Dial(cg.RPCConfig.Address, grpc.WithInsecure())
	if err != nil {
		log.Panic("did not connect.", err)
		return nil
	}
	// defer conn.Close()
	cli := &ChainCli{}

	cli.ChainClient = pd.NewChainClient(conn)
	return cli
}

// Close close chain client
func (cli *ChainCli) Close() error {
	err := cli.conn.Close()
	return err
}
