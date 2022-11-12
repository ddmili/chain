package chain

import "chain/internal/store"

//区块迭代器
type blockchainIterator struct {
	CurrentBlockHash []byte
	BD               *store.BlockchainDB
}

// NewBlockchainIterator 获取区块迭代器实例
func NewBlockchainIterator(bc *Blockchain) *blockchainIterator {
	blockchainIterator := &blockchainIterator{bc.BD.View([]byte(LastBlockHashMapping), store.BlockBucket), bc.BD}
	return blockchainIterator
}

// Next 迭代下一个区块信息
func (bi *blockchainIterator) Next() *Block {
	currentByte := bi.BD.View(bi.CurrentBlockHash, store.BlockBucket)
	if len(currentByte) == 0 {
		return nil
	}
	block := Block{}
	block.Deserialize(currentByte)
	bi.CurrentBlockHash = block.PreHash
	return &block
}
