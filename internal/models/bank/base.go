package bank

import (
	"blockchain-server/internal/blockchain/client"

	"github.com/ethereum/go-ethereum/common"
)

type Bank struct {
	client.BlockchainClient
	contractAddress common.Address
}

var ListBank map[string]*Bank = make(map[string]*Bank)
