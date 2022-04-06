package client

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	ethclient "github.com/ethereum/go-ethereum/ethclient"
	rpc "github.com/ethereum/go-ethereum/rpc"
)

type BlockchainClient struct {
	clientABI abi.ABI
	rpcClient *rpc.Client
	ethClient *ethclient.Client
}
