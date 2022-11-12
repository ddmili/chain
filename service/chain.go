package service

import (
	"chain/chain"
	context "context"
)

// ChainService chain service
type ChainService struct {
	bc *chain.Blockchain
}

// Version version of chain
func (bc *ChainService) Version(ctx context.Context, req *VersionRequest) (*VersionResponse, error) {
	return &VersionResponse{}, nil
}

// NewChainService creates a new ChainService
func NewChainService(bc *chain.Blockchain) *ChainService {
	return &ChainService{bc: bc}
}
