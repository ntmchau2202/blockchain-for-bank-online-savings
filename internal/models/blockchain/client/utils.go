package client

import (
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
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

func BuildTxn(opts *bind.TransactOpts, toAddress common.Address, data []byte) (rawTxn *types.Transaction) {
	if toAddress == common.HexToAddress("0x0") {
		rawTxn = types.NewContractCreation(
			opts.Nonce.Uint64(),
			opts.Value,
			opts.GasLimit,
			opts.GasPrice,
			data)
	} else {
		rawTxn = types.NewTransaction(
			opts.Nonce.Uint64(),
			toAddress,
			opts.Value,
			opts.GasLimit,
			opts.GasPrice,
			data,
		)
	}
	return rawTxn
}

func BuildCall(callOpts *bind.TransactOpts, fromAddress, toAddress common.Address, data []byte) (rawCall ethereum.CallMsg) {
	ptrToAddress := toAddress
	return ethereum.CallMsg{
		From:      fromAddress,
		To:        &ptrToAddress,
		GasPrice:  callOpts.GasPrice,
		Gas:       callOpts.GasLimit,
		GasFeeCap: callOpts.GasFeeCap,
		GasTipCap: callOpts.GasTipCap,
		Value:     new(big.Int).SetUint64(callOpts.GasLimit),
		Data:      data,
	}
}
