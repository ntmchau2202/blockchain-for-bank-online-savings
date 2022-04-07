package bank

import (
	"blockchain-server/internal/blockchain/client"

	"github.com/ethereum/go-ethereum/common"
)

type bank struct {
	client.BlockchainClient
	contractAddress common.Address
	name            common.Address
}

var ListBank map[string]*bank = make(map[string]*bank) // map name to its bank instance

// Note: each bank instance is a client
// that always looking for update on the chain made by the blockchain server on this bank's contract
