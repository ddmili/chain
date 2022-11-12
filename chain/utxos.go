/*
	utxo数据库创建的意义在于,不会每次进行转帐时遍历整个区块链,
	而是去utxo数据库查找未消费的交易输出,这样会大大降低性能问题
*/
package chain

import (
	"bytes"
	"chain/internal/log"
	"chain/internal/store"
	"chain/util"
	"encoding/gob"
	"errors"

	"github.com/boltdb/bolt"
)

type UTXOHandle struct {
	BC *Blockchain
}

//重置UTXO数据库
func (u *UTXOHandle) ResetUTXODataBase() {
	//先查找全部未花费UTXO
	utxosMap := u.BC.findAllUTXOs()
	if utxosMap == nil {
		log.Debug("找不到区块,暂不重置UTXO数据库")
		return
	}
	//删除旧的UTXO数据库
	if u.BC.BD.IsBucketExist(store.UTXOBucket) {
		u.BC.BD.DeleteBucket(store.UTXOBucket)
	}
	//创建并将未花费UTXO循环添加
	for k, v := range utxosMap {
		u.BC.BD.Put([]byte(k), u.serialize(v), store.UTXOBucket)
	}
}

// findUTXOFromAddress 根据地址未消费的utxo
func (u *UTXOHandle) findUTXOFromAddress(address string) []*UTXO {
	publicKeyHash := util.GetPublicKeyHashFromAddress(address)
	utxosSlice := []*UTXO{}
	//获取bolt迭代器，遍历整个UTXO数据库
	//打开数据库
	var DBFileName = "blockchain_" + ListenPort + ".db"
	db, err := bolt.Open(DBFileName, 0600, nil)
	if err != nil {
		log.Panic(" findUTXOFromAddress error: %s", err)
	}

	err = db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte(store.UTXOBucket))
		if b == nil {
			return errors.New("datebase view err: not find bucket ")
		}
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			utxos := u.dserialize(v)
			for _, utxo := range utxos {
				if bytes.Equal(utxo.Vout.PublicKeyHash, publicKeyHash) {
					utxosSlice = append(utxosSlice, utxo)
				}
			}
		}
		return nil
	})
	if err != nil {
		log.Panic(" findUTXOFromAddress error: %s", err)
	}
	//关闭数据库
	err = db.Close()
	if err != nil {
		log.Panic("db close err :", err)
	}
	return utxosSlice
}

// Synchrodata 传入交易信息,将交易里的输出添加进utxo数据库,并剔除输入信息
func (u *UTXOHandle) Synchrodata(tss []Transaction) {
	//先将全部输入插入数据库
	for _, ts := range tss {
		utxos := []*UTXO{}
		for index, vOut := range ts.Vout {
			utxos = append(utxos, &UTXO{ts.TxHash, index, vOut})
		}
		u.BC.BD.Put(ts.TxHash, u.serialize(utxos), store.UTXOBucket)
	}

	//在用输出进行剔除
	for _, ts := range tss {
		for _, vIn := range ts.Vint {
			publicKeyHash := util.GeneratePublicKeyHash(vIn.PublicKey)
			//获取bolt迭代器，遍历整个UTXO数据库
			utxoByte := u.BC.BD.View(vIn.TxHash, store.UTXOBucket)
			if len(utxoByte) == 0 {
				log.Panic("Synchrodata err : do not find utxo")
			}
			utxos := u.dserialize(utxoByte)
			newUTXO := []*UTXO{}
			for _, utxo := range utxos {
				if utxo.Index == vIn.Index && bytes.Equal(utxo.Vout.PublicKeyHash, publicKeyHash) {
					continue
				}
				newUTXO = append(newUTXO, utxo)
			}
			u.BC.BD.Delete(vIn.TxHash, store.UTXOBucket)
			u.BC.BD.Put(vIn.TxHash, u.serialize(newUTXO), store.UTXOBucket)
		}
	}
}

func (u *UTXOHandle) serialize(utxos []*UTXO) []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(utxos)
	if err != nil {
		panic(err)
	}
	return result.Bytes()
}

func (u *UTXOHandle) dserialize(d []byte) []*UTXO {
	var model []*UTXO
	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(&model)
	if err != nil {
		log.Panic("dserialize error: %s", err)
	}
	return model
}
