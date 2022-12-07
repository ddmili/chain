package service

import (
	"chain/chain"
	"chain/config"
	"chain/service/pd"
	context "context"
	"fmt"
)

// ChainService chain service
type ChainService struct {
	bc *chain.Blockchain
	cg *config.Config
}

// Version version of chain
func (s *ChainService) Version(ctx context.Context, req *pd.VersionRequest) (*pd.VersionResponse, error) {
	version := `
	Git Commit Hash: %s
	Build TimeStamp: %s
	GoLang Version: %s
	`
	version = fmt.Sprintf(version, s.cg.GitHash, s.cg.BuildTime, s.cg.GoVersion)
	return &pd.VersionResponse{Version: version}, nil
}

// // GenerateWallet creates a new wallet
// func (s *ChainService) GenerateWallet(ctx context.Context, req *pd.VersionRequest) (*pd.VersionResponse, error) {
// 	address, privkey, mnemonicWord := s.bc.Wallets.GenerateWallet()
// 	fmt.Println("助记词：", mnemonicWord)
// 	fmt.Println("私钥：", privkey)
// 	fmt.Println("地址：", address)
// }

// NewChainService creates a new ChainService
func NewChainService(bc *chain.Blockchain, cg *config.Config) *ChainService {
	return &ChainService{bc: bc, cg: cg}
}
