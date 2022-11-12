package store

import (
	"chain/internal/log"
	"errors"

	"github.com/boltdb/bolt"
)

// Put 存入数据
func (s *BlockchainDB) Put(k, v []byte, bt BucketType) {
	db, err := bolt.Open(s.DBFileName(), 0600, nil)
	defer db.Close()
	if err != nil {
		log.Panic("db error: %s", err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bt))
		if bucket == nil {
			var err error
			bucket, err = tx.CreateBucket([]byte(bt))
			if err != nil {
				log.Panic("db error: %s", err)
			}
		}
		err := bucket.Put(k, v)
		if err != nil {
			log.Panic("db error: %s", err)
		}
		return nil
	})
	if err != nil {
		log.Panic("db error: %s", err)
	}
}

// View 查看数据
func (s *BlockchainDB) View(k []byte, bt BucketType) []byte {
	db, err := bolt.Open(s.DBFileName(), 0600, nil)
	defer db.Close()
	if err != nil {
		log.Panic("db error: %s", err)
	}

	result := []byte{}
	err = db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bt))
		if bucket == nil {
			return errors.New("db error:" + string(bt))
		}
		result = bucket.Get(k)
		return nil
	})
	if err != nil {
		return nil
	}
	//不再次赋值的话，返回值会报错，不知道狗日的啥意思
	realResult := make([]byte, len(result))
	copy(realResult, result)
	return realResult
}

// Delete 删除数据
func (s *BlockchainDB) Delete(k []byte, bt BucketType) bool {
	db, err := bolt.Open(s.DBFileName(), 0600, nil)
	defer db.Close()
	if err != nil {
		log.Panic("db error: %s", err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bt))
		if bucket == nil {
			msg := "datebase delete warnning:没有找到仓库：" + string(bt)
			return errors.New(msg)
		}
		err := bucket.Delete(k)
		if err != nil {
			log.Panic("db error: %s", err)
		}
		return nil
	})
	if err != nil {
		log.Panic("db error: %s", err)
	}
	return true
}

// DeleteBucket 删除仓库
func (s *BlockchainDB) DeleteBucket(bt BucketType) bool {
	db, err := bolt.Open(s.DBFileName(), 0600, nil)
	defer db.Close()
	if err != nil {
		log.Panic("db error: %s", err)
	}
	err = db.Update(func(tx *bolt.Tx) error {
		err := tx.DeleteBucket([]byte(bt))
		if err != nil {
			log.Panic("db error: %s", err)
		}
		return nil
	})
	if err != nil {
		log.Panic("db error: %s", err)
	}

	return true
}
