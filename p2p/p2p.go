package p2p

import (
	"chain/config"
	"context"
	"fmt"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/multiformats/go-multiaddr"
	"go.uber.org/fx"
)

type command string

//网络通讯互相发送的命令
const (
	cVersion     command = "version"
	cGetHash     command = "getHash"
	cHashMap     command = "hashMap"
	cGetBlock    command = "getBlock"
	cBlock       command = "block"
	cTransaction command = "transaction"
	cMyError     command = "myError"
)

//发送数据的头部多少位为命令
const prefixCMDLength = 12

var ProtocolID = "/chain/1.1.0"

// Network p2p network
type Network struct {
	ctx      context.Context
	host     host.Host
	cg       *config.Config
	peerPool map[string]peer.AddrInfo
}

// NewNetwork create a p2p network
func NewNetwork(ctx context.Context, cg *config.Config, lc fx.Lifecycle) *Network {
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

	net := &Network{host: node, ctx: ctx}
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

	net.init()

	return net
}

func (n *Network) init() {

	//发现对等节点
	n.peerPool = make(map[string]peer.AddrInfo)

	go n.findP2PPeer()
}

// Close close network
func (n *Network) Close() error {
	err := n.host.Close()
	if err != nil {
		fmt.Println("p2p shutdown error:", err)
	}
	return err
}

//启动mdns寻找p2p网络 并等节点连接
func (n *Network) findP2PPeer() {
	peerChan := n.mdns(n.ctx, n.host, RendezvousString)
	for {
		peer := <-peerChan // will block untill we discover a peer
		//将发现的节点加入节点池
		n.peerPool[fmt.Sprint(peer.ID)] = peer
	}
}

//默认前十二位为命令名称
func jointMessage(cmd command, content []byte) []byte {
	b := make([]byte, prefixCMDLength)
	for i, v := range []byte(cmd) {
		b[i] = v
	}
	joint := make([]byte, 0)
	joint = append(b, content...)
	return joint
}

//默认前十二位为命令名称
func splitMessage(message []byte) (cmd string, content []byte) {
	cmdBytes := message[:prefixCMDLength]
	newCMDBytes := make([]byte, 0)
	for _, v := range cmdBytes {
		if v != byte(0) {
			newCMDBytes = append(newCMDBytes, v)
		}
	}
	cmd = string(newCMDBytes)
	content = message[prefixCMDLength:]
	return
}
