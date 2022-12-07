package wallet

import (
	"chain/config"
	"chain/internal/log"
	"chain/internal/store"
	"errors"
)

const walletListKey = "wallet_list"

// Wallets 钱包
type Wallets struct {
	bd     *store.BlockchainDB
	bucket store.BucketType
	cg     *config.WalletConfig
}

// NewWallets 创建一个新钱包实例
func NewWallets(bd *store.BlockchainDB, cg *config.Config) *Wallets {
	return &Wallets{bd: bd, bucket: store.BucketType(cg.WalletsBucket), cg: cg.WalletConfig}
}

// GetWallet 获取钱包信息
func (w *Wallets) GetWallet(address string) (*Wallet, error) {
	if w.bd.IsBucketExist(w.bucket) {
		listBytes := w.bd.View([]byte(address), w.bucket)
		if len(listBytes) == 0 {
			return nil, errors.New("wallet not found")
		}
		wallet := &Wallet{}
		wallet.Deserialize(listBytes)
		return wallet, nil
	}
	return nil, errors.New("wallet bucket not exist")
}

// GetAllWallet 获取全部钱包
func (w *Wallets) GetAllWallet() []*Wallet {
	list := []*Wallet{}
	w.bd.AllData(w.bucket, func(data []byte) {
		tmp := &Wallet{}
		tmp.Deserialize(data)
		list = append(list, tmp)
	})

	return list
}

// GenerateWallet 创建公私钥实例
func GenerateWallet(wordPath string) *Wallet {
	b := &Wallet{nil, nil, nil}
	b.MnemonicWord = getChineseMnemonicWord(wordPath)
	b.newKeyPair()
	return b
}

// GenerateWallet 生成钱包
func (w *Wallets) GenerateWallet() (address, privKey, mnemonicWord string) {
	wallet := GenerateWallet(w.cg.WordPath)
	if wallet == nil {
		log.Fatal("创建钱包失败，检查助记词是否符合创建规则！")
	}
	privKey = wallet.GetPrivateKey()
	addressByte := wallet.GetAddress()
	w.storage(addressByte, wallet)
	//将地址存入实例
	address = string(addressByte)
	//将助记词拼接成json格式并返回
	mnemonicWord = "["
	for i, v := range wallet.MnemonicWord {
		mnemonicWord += "\"" + v + "\""
		if i != len(wallet.MnemonicWord)-1 {
			mnemonicWord += ","
		} else {
			mnemonicWord += "]"
		}
	}
	return
}

//将钱包信息存入数据库
func (w *Wallets) storage(address []byte, keys *Wallet) {
	b := w.bd.View(address, w.bucket)
	if len(b) != 0 {
		log.Warn("钱包早已存在于数据库中！")
		return
	}
	//将公私钥以地址为键 存入数据库
	w.bd.Put(address, keys.serliazle(), w.bucket)
}
