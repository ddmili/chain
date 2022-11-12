package config

import (
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

const (
	// DefaultConfig 默认配置文件
	DefaultConfig = "./config.yml"
	// ConfigEnv 默认配置文件路径
	ConfigEnv = "CHAIN_CONFIG"
)

// Config 配置
type Config struct {
	GitHash   string
	BuildTime string
	GoVersion string
	NodeID    int

	*ServerConfig  `yaml:"register"`
	*StoreConfig   `yaml:"datastore"`
	*WalletConfig  `yaml:"wallet"`
	*NetworkConfig `yaml:"network"`
	*RPCConfig     `yaml:"rpc_service"`
}

// ServerConfig 注册
type ServerConfig struct {
	ServiceName string `yaml:"service_name"`
}

// StoreConfig 存储
type StoreConfig struct {
	DbStore string `yaml:"db_store"`
}

// RPCConfig rpc config
type RPCConfig struct {
	Port string `yaml:"port"`
}

// WalletConfig 钱包
type WalletConfig struct {
	WalletsBucket string `yaml:"wallet_bucket"` //钱包表
	WordPath      string `yaml:"word_path"`     //助记词路径
}

// NetworkConfig p2p network
type NetworkConfig struct {
	ListenAddresses string `yaml:"listen_address" `
	BootstrapPeers  string `yaml:"bootstrap_peers" ` //引导节点
}

// NewCg 初始化配置
func NewCg() *Config {
	cgPath := os.Getenv(ConfigEnv)
	if cgPath == "" {
		cgPath = DefaultConfig
	}
	cg, err := NewConfig(cgPath)
	if err != nil {
		panic(err)
	}
	return cg
}

// NewConfig 获取配置
func NewConfig(filepath string) (*Config, error) {
	cfg := &Config{}
	err := cleanenv.ReadConfig(filepath, cfg)
	if err != nil {
		return nil, err
	}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
