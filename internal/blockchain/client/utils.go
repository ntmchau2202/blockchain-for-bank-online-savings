package client

import (
	ethclient "github.com/ethereum/go-ethereum/ethclient"
	rpc "github.com/ethereum/go-ethereum/rpc"
)

func NewClient(url string) (c *BlockchainClient, err error) {
	client, err := rpc.Dial(url)
	if err != nil {
		return
	}
	ethClient := ethclient.NewClient(client)

	c = &BlockchainClient{
		rpcClient: client,
		ethClient: ethClient,
	}
	return
}
