package p2p

import (
	"bufio"
	"chain/internal/log"

	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/protocol"
)

// SendVersionToPeers 向网络中其他节点发送高度信息
func (n *Network) SendVersionToPeers(lastHeight int) {
	newV := version{versionInfo, lastHeight, n.cg.Address}
	data := jointMessage(cVersion, newV.serialize())
	for _, v := range n.peerPool {
		n.SendMessage(v, data)
	}
	log.Trace("version信息发送完毕...")
}

// SendMessage 基础发送信息方法
func (n *Network) SendMessage(peer peer.AddrInfo, data []byte) {
	//连接传入的对等节点
	if err := n.host.Connect(n.ctx, peer); err != nil {
		log.Error("Connection failed:", err)
		return
	}
	//打开一个流，向流写入信息后关闭
	stream, err := n.host.NewStream(n.ctx, peer.ID, protocol.ID(ProtocolID))
	if err != nil {
		log.Debug("Stream open failed", err)
	} else {
		cmd, _ := splitMessage(data)
		//创建一个缓冲流的容器
		rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))
		//写入信息到缓冲容器
		_, err := rw.Write(data)
		if err != nil {
			log.Error("send message failed:%s", err)
			return
		}
		//向流中写入所有缓冲数据
		err = rw.Flush()
		if err != nil {
			log.Error("send message failed:", err)
			return
		}
		//关闭流，完成一次信息的发送
		err = stream.Close()
		if err != nil {
			log.Error("send message failed:", err)
			return
		}
		log.Debug("send cmd:%s to peer:%v", cmd, peer)
	}
}
