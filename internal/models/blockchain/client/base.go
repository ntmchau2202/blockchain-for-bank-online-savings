package client

import (
	"blockchain-server/internal/models/blockchain/signer"
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
	// get the nonce

	nonce, err := c.ethClient.PendingNonceAt(context.Background(), c.txnSigner.GetOwner())
	if err != nil {
		return
	}
	options.Nonce = new(big.Int).SetUint64(nonce)

	// suggest gas price
	gasPrice, err := c.ethClient.SuggestGasPrice(context.Background())
	if err != nil {
		gasPrice = big.NewInt(100)
		err = nil
	}
	options.GasPrice = gasPrice

	// suggest gas tip cap
	gasTipCap, err := c.ethClient.SuggestGasTipCap(context.Background())
	if err != nil {
		gasTipCap = big.NewInt(0)
		err = nil
	}
	options.GasTipCap = gasTipCap

	// set timeout
	timeOut, cancelFunc := context.WithTimeout(context.Background(), timeOutThreshold)
	options.Context = timeOut // 5s timeout
	options.NoSend = false
	return
}

func (c *BlockchainClient) EstimateGasLimit(toAddress string, dataToBeSent []byte) (gasLimit *big.Int, err error) {
	if c.txnSigner == nil {
		return nil, errors.New("no signer asscociated")
	}
	toAddr := common.HexToAddress(toAddress)
	msg := ethereum.CallMsg{
		From:     c.txnSigner.GetOwner(),
		To:       &toAddr,
		GasPrice: c.txnSigner.GetTxnOptsOfSigner().GasPrice,
		Data:     dataToBeSent,
	}
	val, err := c.ethClient.EstimateGas(context.Background(), msg)
	if err != nil {
		val = 21000
		txNonZero := 68
		additionalGas := new(big.Int).SetInt64(int64(txNonZero))
		additionalGas.Mul(additionalGas, big.NewInt(int64(len(dataToBeSent))))
		gasLimit = new(big.Int).SetUint64(val).Add(gasLimit, additionalGas)
		err = nil
		return
	}
	return new(big.Int).SetUint64(val), nil
}

// call only retrieve data
func (c *BlockchainClient) Call(
	method string,
	toAddress string,
	arguments ...interface{},
) (unpacked map[string]interface{}, err error) {
	data, err := c.clientABI.Pack(method, arguments...)
	if err != nil {
		return
	}

	callOpts, _, _ := c.getTxnOpts()
	txn := BuildCall(callOpts, c.txnSigner.GetOwner(), common.HexToAddress(toAddress), data)

	result, err := c.ethClient.CallContract(context.Background(), txn, nil)
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
	toAddress string,
	arguments ...interface{},
) (unpacked map[string]interface{}, err error) {
	data, err := c.clientABI.Pack(method, arguments...)
	if err != nil {
		return
	}

	// sign data before sending
	ptrToAddress := common.HexToAddress(toAddress)
	callOpts, _, _ := c.getTxnOpts()

	txn := BuildTxn(callOpts, ptrToAddress, data)
	signedTxn, err := c.txnSigner.Sign(txn)
	if err != nil {
		return
	}

	err = c.ethClient.SendTransaction(context.Background(), signedTxn)
	if err != nil {
		return
	}

	receipt, err := c.ethClient.TransactionReceipt(context.Background(), txn.Hash())
	if err != nil {
		return
	}

	unpacked = make(map[string]interface{})
	receivedData, err := receipt.MarshalBinary()
	if err != nil {
		return
	}

	unpacked["transaction_id"] = txn.Hash()
	unpacked["block_number"] = receipt.BlockNumber
	unpacked["post_state"] = receipt.PostState
	unpacked["data"] = receivedData
	return
}
