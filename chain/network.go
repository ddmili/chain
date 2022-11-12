package chain

// Network 用于network包向对等节点发送信息
type Network interface {
	SendVersionToPeers(height int)
	SendTransToPeers(tss []Transaction)
}
