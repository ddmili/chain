package chain

import (
	"chain/internal/log"
	"chain/internal/store"
	"chain/internal/wallet"
	"chain/p2p"
	"chain/util"
	"fmt"
)

// Blockchain 链
type Blockchain struct {
	BD      *store.BlockchainDB //封装的blot结构体
	Network *p2p.Network
	Wallets *wallet.Wallets
}

// NewBlockchain 创建区块链实例
func NewBlockchain(db *store.BlockchainDB, network *p2p.Network, wallets *wallet.Wallets) *Blockchain {
	blockchain := Blockchain{}
	blockchain.BD = db
	blockchain.Network = network
	blockchain.Wallets = wallets
	return &blockchain
}

// CreateGenesisTransaction 创建创世区块交易信息
func (bc *Blockchain) CreateGenesisTransaction(address string, value int) {
	//判断地址格式是否正确
	if !util.IsVailBitcoinAddress(address) {
		log.Error("地址格式不正确:%s\n", address)
		return
	}
	//创世区块数据
	txi := TXInput{[]byte{}, -1, nil, nil}
	//本地一定要存创世区块地址的公私钥信息
	genesisKeys, err := bc.Wallets.GetWallet(address)
	if err != nil {
		log.Fatal("没有找到地址对应的公私钥信息")
		return
	}
	//通过地址获得rip160(sha256(publicKey))
	publicKeyHash := util.GeneratePublicKeyHash(genesisKeys.PublicKey)
	txo := TXOutput{value, publicKeyHash}
	ts := Transaction{nil, []TXInput{txi}, []TXOutput{txo}}
	ts.hash()
	tss := []Transaction{ts}
	//开始生成区块链的第一个区块
	bc.newGenesisBlockchain(tss)
	//创世区块后,更新本地最新区块为1并,向全网节点发送当前区块链高度1
	NewestBlockHeight = 1
	bc.Network.SendVersionToPeers(1)
	fmt.Println("已成生成创世区块")
	//重置utxo数据库，将创世数据存入
	utxos := UTXOHandle{bc}
	utxos.ResetUTXODataBase()
}

// newGenesisBlockchain 创建区块链
func (bc *Blockchain) newGenesisBlockchain(transaction []Transaction) {
	//判断一下是否已生成创世区块
	if len(bc.BD.View([]byte(LastBlockHashMapping), store.BlockBucket)) != 0 {
		log.Fatal("不可重复生成创世区块")
	}
	//生成创世区块
	genesisBlock := newGenesisBlock(transaction)
	//添加到数据库中
	bc.AddBlock(genesisBlock)
}

// AddBlock 添加区块信息到数据库，并更新lastHash
func (bc *Blockchain) AddBlock(block *Block) {
	bc.BD.Put(block.Hash, block.Serialize(), store.BlockBucket)
	bci := NewBlockchainIterator(bc)
	currentBlock := bci.Next()
	if currentBlock == nil || currentBlock.Height < block.Height {
		bc.BD.Put([]byte(LastBlockHashMapping), block.Hash, store.BlockBucket)
	}
}

// findAllUTXOs 查找数据库中全部未花费的UTXO
func (bc *Blockchain) findAllUTXOs() map[string][]*UTXO {
	utxosMap := make(map[string][]*UTXO)
	txInputmap := make(map[string][]TXInput)
	bcIterator := NewBlockchainIterator(bc)
	for {
		currentBlock := bcIterator.Next()
		if currentBlock == nil {
			return nil
		}
		//必须倒序 否则有的已花费不会被扣掉
		for i := len(currentBlock.Transactions) - 1; i >= 0; i-- {
			var utxos = []*UTXO{}
			ts := currentBlock.Transactions[i]
			for _, vInt := range ts.Vint {
				txInputmap[string(vInt.TxHash)] = append(txInputmap[string(vInt.TxHash)], vInt)
			}

		VoutTag:
			for index, vOut := range ts.Vout {
				if txInputmap[string(ts.TxHash)] == nil {
					utxos = append(utxos, &UTXO{ts.TxHash, index, vOut})
				} else {
					for _, vIn := range txInputmap[string(ts.TxHash)] {
						if vIn.Index == index {
							continue VoutTag
						}
					}
					utxos = append(utxos, &UTXO{ts.TxHash, index, vOut})
				}
				utxosMap[string(ts.TxHash)] = utxos
			}
		}

		if isGenesisBlock(currentBlock) {
			break
		}
	}
	return utxosMap
}
