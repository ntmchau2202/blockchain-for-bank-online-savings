package client

import (
	"blockchain-server/internal/blockchain/signer"
	"context"
	"errors"
	"math/big"
	"os"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethclient "github.com/ethereum/go-ethereum/ethclient"
	rpc "github.com/ethereum/go-ethereum/rpc"
)

type BlockchainClient struct {
	rpcClient *rpc.Client
	ethClient *ethclient.Client
	clientABI abi.ABI
	txnSigner *signer.Signer
}

const timeOutThreshold = 5 * time.Second

func (c *BlockchainClient) AddSignerWithChainID(privateKeyPath string, chainId int64) (err error) {
	c.txnSigner, err = signer.NewSignerWithChainID(privateKeyPath, chainId)
	if err != nil {
		return
	}
	return
}

func (c *BlockchainClient) AddABI(path string) (err error) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	parsed, err := abi.JSON(file)
	if err != nil {
		return
	}
	c.clientABI = parsed
	return nil
}

func (c *BlockchainClient) getTxnOpts() (options *bind.TransactOpts, cancelFunc context.CancelFunc, err error) {
	if c.txnSigner == nil {
		return nil, nil, errors.New("no signer associated")
	}

	options = c.txnSigner.GetTxnOptsOfSigner()
	timeOut, cancelFunc := context.WithTimeout(context.Background(), timeOutThreshold)
	options.NoSend = false
	options.Context = timeOut // 5s timeout
	return
}

// call only retrieve data
func (c *BlockchainClient) Call(
	method string,
	fromAddress string,
	toAddress string,
	arguments ...interface{},
) (unpacked map[string]interface{}, err error) {
	data, err := c.clientABI.Pack(method, arguments...)
	if err != nil {
		return
	}

	// get gas price
	gasPrice, err := c.ethClient.SuggestGasPrice(context.Background())
	if err != nil {
		return
	}
	// get gas limit
	gasLimit := new(big.Int).SetInt64(68)
	gasLimit = gasLimit.Mul(gasLimit, new(big.Int).SetInt64(int64(len(data))))

	// convert to ptr
	ptrToAddress := common.HexToAddress(toAddress)
	callMsg := ethereum.CallMsg{
		From:      common.HexToAddress(fromAddress),
		To:        &ptrToAddress,
		GasPrice:  gasPrice,
		Gas:       gasLimit.Uint64(),
		GasFeeCap: big.NewInt(500000),
		GasTipCap: big.NewInt(0),
		Value:     gasLimit,
		Data:      data,
	}

	result, err := c.ethClient.CallContract(context.Background(), callMsg, nil)
	if err != nil {
		return
	}
	unpacked = make(map[string]interface{})
	err = c.clientABI.UnpackIntoMap(unpacked, method, result)
	return
}

// send will perform the transaction
func (c *BlockchainClient) Send(
	method string,
	fromAddress string,
	toAddress string,
	arguments ...interface{},
) (unpacked map[string]interface{}, err error) {
	data, err := c.clientABI.Pack(method, arguments...)
	if err != nil {
		return
	}

	// get gas price
	gasPrice, err := c.ethClient.SuggestGasPrice(context.Background())
	if err != nil {
		return
	}
	// get gas limit
	gasLimit := new(big.Int).SetInt64(68)
	gasLimit = gasLimit.Mul(gasLimit, new(big.Int).SetInt64(int64(len(data))))

	// convert to ptr
	ptrToAddress := common.HexToAddress(toAddress)
	callMsg := ethereum.CallMsg{
		From:      common.HexToAddress(fromAddress),
		To:        &ptrToAddress,
		GasPrice:  gasPrice,
		Gas:       gasLimit.Uint64(),
		GasFeeCap: big.NewInt(500000),
		GasTipCap: big.NewInt(0),
		Value:     gasLimit,
		Data:      data,
	}

	result, err := c.ethClient.CallContract(context.Background(), callMsg, nil)
	if err != nil {
		return
	}
	unpacked = make(map[string]interface{})
	err = c.clientABI.UnpackIntoMap(unpacked, method, result)
	return
}
