package chain

import (
	"math"
)

// NewestBlockHeight 当前节点发现的网络中最新区块高度
var NewestBlockHeight int

// ListenPort 当前本地监听端口
var ListenPort string

// TokenRewardNum 挖矿奖励代币数量
var TokenRewardNum int

// TargetBits 挖矿难度值
var TargetBits uint

// ChineseMnwordPath 中文助记词地址
var ChineseMnwordPath string

// RewardAddrMapping 奖励地址在数据库中的键
const RewardAddrMapping = "rewardAddress"

// LastBlockHashMapping 最新区块Hash在数据库中的键
const LastBlockHashMapping = "lastHash"

// addrListMapping 钱包地址在数据库中的键
const addrListMapping = "addressList"

// version 公链版本信息默认为0
const version = byte(0x00)

// checkSum 两次sha256(公钥hash)后截取的字节数量
const checkSum = 4

// maxInt 随机数不能超过的最大值
const maxInt = math.MaxInt64
