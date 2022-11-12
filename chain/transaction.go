package chain

import (
	"bytes"
	"chain/internal/log"
	"chain/util"
	"crypto/sha256"
	"encoding/gob"
)

// Transaction 交易列表信息
type Transaction struct {
	TxHash []byte
	// Vint UTXO输入
	Vint []TXInput
	// Vout UTXO输出
	Vout []TXOutput
}

// hash 对此笔交易的输入,输出进行hash运算后存入交易hash(txhash)
func (t *Transaction) hash() {
	tBytes := t.Serialize()
	//加入随机数byte
	randomNumber := util.GenerateRealRandom()
	randomByte := util.Int64ToBytes(randomNumber)
	sumByte := bytes.Join([][]byte{tBytes, randomByte}, []byte(""))
	hashByte := sha256.Sum256(sumByte)
	t.TxHash = hashByte[:]
}

// hashSign 作为数字签名的hash方法，为什么不用gob序列化后hash，因为涉及到tcp传输gob直接序列化有问题，所以单独拼接成byte数组再hash
func (t *Transaction) hashSign() []byte {
	t.TxHash = nil
	nHash := []byte{}
	for _, v := range t.Vint {
		nHash = append(nHash, v.TxHash...)
		nHash = append(nHash, v.PublicKey...)
		nHash = append(nHash, util.Int64ToBytes(int64(v.Index))...)
	}
	for _, v := range t.Vout {
		nHash = append(nHash, v.PublicKeyHash...)
		nHash = append(nHash, util.Int64ToBytes(int64(v.Value))...)
	}
	hashByte := sha256.Sum256(nHash)
	return hashByte[:]
}

// Serialize 将transaction序列化成[]byte
func (t *Transaction) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(t)
	if err != nil {
		panic(err)
	}
	return result.Bytes()
}

// getTransBytes 将整笔交易里的成员依次转换成字节数组,拼接成整体后 返回
func (t *Transaction) getTransBytes() []byte {
	if t.TxHash == nil || t.Vout == nil {
		log.Panic("交易信息不完整，无法拼接成字节数组")
		return nil
	}
	transBytes := []byte{}
	transBytes = append(transBytes, t.TxHash...)
	for _, v := range t.Vint {
		transBytes = append(transBytes, v.TxHash...)
		transBytes = append(transBytes, util.Int64ToBytes(int64(v.Index))...)
		transBytes = append(transBytes, v.Signature...)
		transBytes = append(transBytes, v.PublicKey...)
	}
	for _, v := range t.Vout {
		transBytes = append(transBytes, util.Int64ToBytes(int64(v.Value))...)
		transBytes = append(transBytes, v.PublicKeyHash...)
	}
	return transBytes
}

// customCopy 从原交易里拷贝出一个新的交易
func (t *Transaction) customCopy() Transaction {
	newVin := []TXInput{}
	newVout := []TXOutput{}
	for _, vin := range t.Vint {
		newVin = append(newVin, TXInput{vin.TxHash, vin.Index, nil, nil})
	}
	for _, vout := range t.Vout {
		newVout = append(newVout, TXOutput{vout.Value, vout.PublicKeyHash})
	}
	return Transaction{t.TxHash, newVin, newVout}
}

// isGenesisTransaction 判断是否是创世区块的交易
func isGenesisTransaction(tss []Transaction) bool {
	if tss != nil {
		if tss[0].Vint[0].Index == -1 {
			return true
		}
	}
	return false
}
