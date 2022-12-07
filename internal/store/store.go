package store

import (
	"chain/config"
	"chain/internal/log"
	"fmt"
	"os"

	"github.com/boltdb/bolt"
)

// BucketType 仓库类型
type BucketType string

const (
	BlockBucket BucketType = "blocks"
	AddrBucket  BucketType = "address"
	UTXOBucket  BucketType = "utxo"
)

// BlockchainDB db
type BlockchainDB struct {
	conf   *config.StoreConfig
	nodeID int
}

// NewDb create db
func NewDb(cg *config.Config) *BlockchainDB {
	return &BlockchainDB{conf: cg.StoreConfig, nodeID: cg.NodeID}
}

// DBFileName db name
func (s *BlockchainDB) DBFileName() string {
	return fmt.Sprintf("%s/chain_%d.db", s.conf.DbStore, s.nodeID)
}

// IsBlotExist 判断数据库是否存在
func (s BlockchainDB) IsBlotExist(nodeID string) bool {
	_, err := os.Stat(s.DBFileName())
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

// IsBucketExist 判断仓库是否存在
func (s BlockchainDB) IsBucketExist(bt BucketType) bool {
	var isBucketExist bool

	db, err := bolt.Open(s.DBFileName(), 0600, nil)
	if err != nil {
		log.Panic("db error:%s", err)
	}

	err = db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bt))
		if bucket == nil {
			isBucketExist = false
		} else {
			isBucketExist = true
		}
		return nil
	})
	if err != nil {
		log.Panic("db IsBucketExist err:%s", err)
	}

	err = db.Close()
	if err != nil {
		log.Panic("db close err :%s", err)
	}
	return isBucketExist
}
