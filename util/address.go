package util

import (
	"bytes"
	"chain/constant"
	"crypto/sha256"
)

// IsVailBitcoinAddress 判断是否是有效的比特币地址
func IsVailBitcoinAddress(address string) bool {
	addressByte := []byte(address)
	fullHash := Base58Decode(addressByte)
	if len(fullHash) != 25 {
		return false
	}
	prefixHash := fullHash[:len(fullHash)-constant.CheckSum]
	tailHash := fullHash[len(fullHash)-constant.CheckSum:]
	tailHash2 := CheckSumHash(prefixHash)
	if bytes.Compare(tailHash, tailHash2[:]) == 0 {
		return true
	} else {
		return false
	}
}

// CheckSumHash 检测
func CheckSumHash(versionPublicKeyHash []byte) []byte {
	versionPublicKeyHashSha1 := sha256.Sum256(versionPublicKeyHash)
	versionPublicKeyHashSha2 := sha256.Sum256(versionPublicKeyHashSha1[:])
	tailHash := versionPublicKeyHashSha2[:constant.CheckSum]
	return tailHash
}

// GeneratePublicKeyHash 公钥
func GeneratePublicKeyHash(publicKey []byte) []byte {
	sha256PubKey := sha256.Sum256(publicKey)
	r := NewRipemd160()
	r.Reset()
	r.Write(sha256PubKey[:])
	ripPubKey := r.Sum(nil)
	return ripPubKey
}

// GetPublicKeyHashFromAddress 获取公钥
func GetPublicKeyHashFromAddress(address string) []byte {
	addressBytes := []byte(address)
	fullHash := Base58Decode(addressBytes)
	publicKeyHash := fullHash[1 : len(fullHash)-constant.CheckSum]
	return publicKeyHash
}
