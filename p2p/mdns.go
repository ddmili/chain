package p2p

import (
	"context"

	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
)

var (
	RendezvousString = "meetme"
)

type discoveryNotifee struct {
	PeerChan chan peer.AddrInfo
}

//当网络中找到新节点时,此方法会被调用
func (n *discoveryNotifee) HandlePeerFound(pi peer.AddrInfo) {
	n.PeerChan <- pi
}

// mdns 开启mdns
func (n *Network) mdns(ctx context.Context, peerhost host.Host, rendezvous string) chan peer.AddrInfo {
	//注册Notifee接口类型
	d := &discoveryNotifee{}

	d.PeerChan = make(chan peer.AddrInfo)
	// time.Second检索当前网络节点的频率
	mdns.NewMdnsService(peerhost, rendezvous, d)

	return d.PeerChan
}
