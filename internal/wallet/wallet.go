package wallet

import (
	"bytes"
	"chain/constant"
	"chain/internal/log"
	"chain/util"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/gob"
	"math/big"
	"os"
)

// Wallet 钱包
type Wallet struct {
	PrivateKey   *ecdsa.PrivateKey
	PublicKey    []byte
	MnemonicWord []string
}

// CreateBitcoinKeysByMnemonicWord 根据助记词创建公私钥
func CreateBitcoinKeysByMnemonicWord(mnemonicWord []string) *Wallet {
	if len(mnemonicWord) != 7 {
		log.Error("助记词格式不正确，应为七对中文双字词语")
		return nil
	}
	for _, v := range mnemonicWord {
		if len(v) != 6 {
			log.Error("助记词格式不正确，应为七对中文双字词语")
			return nil
		}
	}

	b := &Wallet{nil, nil, nil}
	b.MnemonicWord = mnemonicWord
	b.newKeyPair()
	return b
}

// newKeyPair 根据中文助记词生成公私钥对
func (b *Wallet) newKeyPair() {
	curve := elliptic.P256()
	var err error
	buf := bytes.NewReader(b.jointSpeed())
	b.PrivateKey, err = ecdsa.GenerateKey(curve, buf)
	if err != nil {
		log.Panic("newKeyPair error: %s")
	}
	b.PublicKey = append(b.PrivateKey.PublicKey.X.Bytes(), b.PrivateKey.PublicKey.Y.Bytes()...)
}

// jointSpeed 将助记词拼接成字节数组，并截取前40位
func (b Wallet) jointSpeed() []byte {
	bs := make([]byte, 0)
	for _, v := range b.MnemonicWord {
		bs = append(bs, []byte(v)...)
	}
	return bs[:40]
}

//私钥长度为32字节
const privKeyBytesLen = 32

// GetPrivateKey 获取私钥
func (w *Wallet) GetPrivateKey() string {
	d := w.PrivateKey.D.Bytes()
	b := make([]byte, 0, privKeyBytesLen)
	priKey := paddedAppend(privKeyBytesLen, b, d)
	//base58加密
	return string(util.Base58Encode(priKey))
}

// getAddress 通过公钥获得地址
func (b *Wallet) getAddress() []byte {
	//1.ripemd160(sha256(publicKey))
	ripPubKey := util.GeneratePublicKeyHash(b.PublicKey)
	//2.最前面添加一个字节的版本信息获得 versionPublicKeyHash
	versionPublicKeyHash := append([]byte{constant.Version}, ripPubKey[:]...)
	//3.sha256(sha256(versionPublicKeyHash))  取最后四个字节的值
	tailHash := util.CheckSumHash(versionPublicKeyHash)
	//4.拼接最终hash versionPublicKeyHash + checksumHash
	finalHash := append(versionPublicKeyHash, tailHash...)
	//进行base58加密
	address := util.Base58Encode(finalHash)
	return address
}

// serliazle 序列化
func (b *Wallet) serliazle() []byte {
	var result bytes.Buffer
	gob.Register(elliptic.P256())
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(b)
	if err != nil {
		panic(err)
	}
	return result.Bytes()
}

// Deserialize 反序列化
func (v *Wallet) Deserialize(d []byte) {
	decoder := gob.NewDecoder(bytes.NewReader(d))
	gob.Register(elliptic.P256())
	err := decoder.Decode(v)
	if err != nil {
		log.Panic("Deserialize error: %s", err)
	}
}

// GetAddressFromPublicKey 通过公钥信息获得地址
func GetAddressFromPublicKey(publicKey []byte) string {
	if publicKey == nil {
		return ""
	}
	b := Wallet{PublicKey: publicKey}
	return string(b.getAddress())
}

func paddedAppend(size uint, dst, src []byte) []byte {
	for i := 0; i < int(size)-len(src); i++ {
		dst = append(dst, 0)
	}
	return append(dst, src...)
}

// getChineseMnemonicWord 获取中文种子
func getChineseMnemonicWord(chineseWordPath string) []string {
	file, err := os.Open(chineseWordPath)
	//file,err:=os.Open("D:/programming/golang/GOPATH/src/github.com/corgi-kx/blockchain_golang/blc/chinese_mnemonic_world.txt")
	if err != nil {
		log.Panic("getChineseMnemonicWord error: %s", err)
	}
	s := []string{}
	//因为种子最高40个字节，所以就取7对词语，7*2*3 = 42字节，返回后在截取前40位
	for i := 0; i < 7; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(5948)) //词库一共5949对词语，顾此设置随机数最高5948
		if err != nil {
			log.Panic("getChineseMnemonicWord error: %s")
		}
		b := make([]byte, 6)
		_, err = file.ReadAt(b, n.Int64()*7+3) //从文件的具体位置读取 防止乱码
		if err != nil {
			log.Panic("getChineseMnemonicWord error: %s")
		}
		s = append(s, string(b))
	}
	file.Close()
	return s
}
