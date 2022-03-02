package client

import (
	"context"
	"errors"
	"fmt"
	"gr-blockchain-side/internal/blockchain/signer"
	"gr-blockchain-side/internal/config"
	"log"
	"math/big"

	"os"
	"time"

	ether "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	ethclient "github.com/ethereum/go-ethereum/ethclient"
	rpc "github.com/ethereum/go-ethereum/rpc"
)

type BlockchainClient struct {
	client    *rpc.Client
	ethClient *ethclient.Client
	signer    *signer.Signer
	abi       abi.ABI
	contract  ethereum.Address
}

func NewBlockchainClient() (c *BlockchainClient, err error) {
	client, err := rpc.Dial(config.DefaultConfig.Node)
	if err != nil {
		return
	}
	ethClient := ethclient.NewClient(client)
	file, err := os.Open(config.DefaultConfig.ABIPath)
	if err != nil {
		return
	}
	parsed, err := abi.JSON(file)
	if err != nil {
		return
	}

	c = &BlockchainClient{
		client:    client,
		ethClient: ethClient,
		signer:    config.DefaultConfig.Signer,
		abi:       parsed,
		contract:  ethereum.HexToAddress(config.DefaultConfig.Contract),
	}
	return
}

func (bc BlockchainClient) GetRPCClient() (c *rpc.Client) {
	return bc.client
}

func (bc BlockchainClient) GetETHClent() (c *ethclient.Client) {
	return bc.ethClient
}

func (bc BlockchainClient) SaveSavingsAccountCreation(message AddSavingsAccountMessage) (trueResult *types.Transaction, err error) {
	opts, cancel, err := bc.getTxnOpts()
	defer cancel()

	if err != nil {
		return
	} else {
		tx, err := bc.buildTx(
			opts,
			"registerSavingsAccount",
			message.Details.BankAccountNumber,
			message.Details.SavingsAccountNumber,
			time.Now().String(), // time issued
			message.Details.TimeCreated,
			message.Details.InitialAmount,
			message.Details.SavingsPeriod,
			message.Details.InterestRate,
			message.Details.TypeOfSavings,
			message.Details.TransactionUnit,
		)
		if err != nil {
			return nil, err
		} else {
			return bc.signAndBroadcast(tx, bc.signer)
		}
	}
	return
}

func (bc BlockchainClient) SaveSavingsAccountSettlement(message SettleSavingsAccountMessage) (trueResult *types.Transaction, err error) {
	opts, cancel, err := bc.getTxnOpts()
	defer cancel()

	if err != nil {
		return
	} else {

		tx, err := bc.buildTx(
			opts,
			"settleSavingsAccount",
			message.Details.BankAccountNumber,
			message.Details.SavingsAccountNumber,
			time.Now().String(), // time issued,
			message.Details.TimeSettled,
			message.Details.InterestAmount,
			message.Details.TotalAmount,
		)
		if err != nil {
			return nil, err
		} else {
			return bc.signAndBroadcast(tx, bc.signer)
		}
	}
	return
}

// func (bc BlockchainClient) GetTransactionsOnBankAccount(bankAccount string) (trueResult []interface{}, err error) {
// 	input, err := bc.abi.Pack("getTransactionsOnBankAccount", bankAccount)
// 	if err != nil {
// 		return
// 	}
// 	value := big.NewInt(0)
// 	from := bc.signer.GetAddress()
// 	msg := ether.CallMsg{
// 		From:  from,
// 		To:    &bc.contract,
// 		Value: value,
// 		Data:  input,
// 	}
// 	result, err := bc.ethClient.CallContract(context.Background(), msg, nil)
// 	if err != nil {
// 		return
// 	}

// 	trueResult, err = bc.abi.Unpack("getTransactionsOnBankAccount", result)
// 	return
// }

// func (bc BlockchainClient) GetTransactionsOnSavingsAccount(savingAccount string /*args here*/) (trueResult []interface{}, err error) {
// 	input, err := bc.abi.Pack("getTransactionsOnSavingsAccount", savingAccount)
// 	if err != nil {
// 		return
// 	}
// 	value := big.NewInt(0)
// 	from := bc.signer.GetAddress()
// 	msg := ether.CallMsg{
// 		From:  from,
// 		To:    &bc.contract,
// 		Value: value,
// 		Data:  input,
// 	}
// 	result, err := bc.ethClient.CallContract(context.Background(), msg, nil)
// 	if err != nil {
// 		return
// 	}

// 	trueResult, err = bc.abi.Unpack("getTransactionsOnSavingsAccount", result)
// 	return
// }

func (bc BlockchainClient) getTxnOpts() (result *bind.TransactOpts, cancel context.CancelFunc, err error) {
	shared := bc.signer.GetTransactOpts()
	nonce, err := bc.getNonceFromNode()
	if err != nil {
		return nil, func() {}, err
	}
	gasLimit := uint64(shared.GasLimit)
	gasFeeCap := big.NewInt(5000000000)
	gasPrice := big.NewInt(1000000000)
	gasTipCap := big.NewInt(1000000000)
	timeout, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	result = &bind.TransactOpts{
		shared.From,
		nonce,
		shared.Signer,
		shared.Value,
		gasPrice,
		gasFeeCap,
		gasTipCap,
		gasLimit,
		timeout,
		false,
	}
	return result, cancel, nil
}

func (bc *BlockchainClient) getNonceFromNode() (*big.Int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	nonce, err := bc.ethClient.PendingNonceAt(ctx, bc.signer.GetAddress())
	return big.NewInt(int64(nonce)), err
}

func (bc *BlockchainClient) buildTx(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	input, err := bc.abi.Pack(method, params...)
	if err != nil {
		return nil, err
	}
	log.Println("build tx: ", opts.From.Hex())
	return bc.transactTx(opts, &bc.contract, input)
}

func (self *BlockchainClient) transactTx(opts *bind.TransactOpts, contract *ethereum.Address, input []byte) (*types.Transaction, error) {
	var err error
	// Ensure a valid value field and resolve the account nonce
	value := opts.Value
	if value == nil {
		value = new(big.Int)
	}
	var nonce uint64
	if opts.Nonce == nil {
		return nil, errors.New("nonce must be specified")
	} else {
		nonce = opts.Nonce.Uint64()
	}
	// Figure out the gas allowance and gas price values
	gasPrice := opts.GasPrice
	if gasPrice == nil {
		return nil, errors.New("gas price must be specified")
	}
	gasLimit := opts.GasLimit
	if gasLimit == 0 {
		// Gas estimation cannot succeed without code for method invocations
		if contract != nil {
			if code, err := self.ethClient.PendingCodeAt(ensureContext(opts.Context), self.contract); err != nil {
				return nil, err
			} else if len(code) == 0 {
				return nil, bind.ErrNoCode
			}
		}
		// If the contract surely has code (or code is not needed), estimate the transaction
		msg := ether.CallMsg{From: opts.From, To: contract, Value: value, Data: input}
		gasLimit, err = self.ethClient.EstimateGas(ensureContext(opts.Context), msg)
		if err != nil {
			return nil, fmt.Errorf("failed to estimate gas needed: %v", err)
		}
		// add gas limit by 50K gas
		// gasLimit.Add(gasLimit, big.NewInt(50000))
		gasLimit = gasLimit + 50000
	}
	// Create the transaction, sign it and schedule it for execution
	var rawTx *types.Transaction
	if contract == nil {
		rawTx = types.NewContractCreation(nonce, value, gasLimit, gasPrice, input)
	} else {
		rawTx = types.NewTransaction(nonce, self.contract, value, gasLimit, gasPrice, input)
	}
	return rawTx, nil
}

func ensureContext(ctx context.Context) context.Context {
	if ctx == nil {
		return context.TODO()
	}
	return ctx
}

func (self *BlockchainClient) signAndBroadcast(tx *types.Transaction, singer *signer.Signer) (*types.Transaction, error) {
	if tx == nil {
		return nil, errors.New("Nil tx is forbidden here")
	} else {
		signedTx, err := singer.Sign(tx)
		if err != nil {
			return nil, err
		}
		// log.Println("raw tx: ", signedTx)
		ctx := context.Background()
		err = self.ethClient.SendTransaction(ctx, signedTx)
		if err != nil {
			log.Println("failed to broadcast tx: ", err)
		}
		log.Println("send done!")
		return signedTx, nil
	}
}

func parseEvt(contractABI abi.ABI, data []byte) (result interface{}, err error) {
	result = make(map[string]interface{})
	tmp := make(map[string]interface{})
	err = contractABI.UnpackIntoMap(tmp, "OpenSavingsAccount", data)
	if err == nil {
		result = tmp["_openSavingsAccTxn"].(interface{})
		fmt.Println(result, "===============")
		return
	}
	return
}

func (bc BlockchainClient) QueryTxnsByBankAccount(bankAcc string) (fullResult []interface{}, err error) {

	ethClient := bc.GetETHClent()
	bankAccCrypted := ethereum.BytesToHash(crypto.Keccak256([]byte(bankAcc)))
	contractAddress := ethereum.HexToAddress(config.DefaultConfig.Contract)
	query := ether.FilterQuery{
		FromBlock: nil,
		Addresses: []ethereum.Address{contractAddress},
		Topics:    [][]ethereum.Hash{nil, {bankAccCrypted}},
	}
	logs, err := ethClient.FilterLogs(context.Background(), query)
	if err != nil {
		return
	}
	// get ABI to parse
	file, err := os.Open(config.DefaultConfig.ABIPath)
	if err != nil {
		return
	}
	contractABI, err := abi.JSON(file)
	if err != nil {
		return
	}

	for _, record := range logs {
		res, err := parseEvt(contractABI, record.Data)
		if err != nil {
			log.Panic(err)
			continue
		}
		// do something to compare here
		fullResult = append(fullResult, res)
	}
	return fullResult, nil
}

func (bc BlockchainClient) QueryTxnsBySavingsAccount(savingsAcc string) (fullResult []interface{}, err error) {

	ethClient := bc.GetETHClent()
	savingsAccCrypted := ethereum.BytesToHash(crypto.Keccak256([]byte(savingsAcc)))

	// build filter for open savings acc

	contractAddress := ethereum.HexToAddress(config.DefaultConfig.Contract)
	query := ether.FilterQuery{
		FromBlock: nil,
		Addresses: []ethereum.Address{contractAddress},
		Topics:    [][]ethereum.Hash{nil, nil, {savingsAccCrypted}},
	}
	logs, err := ethClient.FilterLogs(context.Background(), query)
	if err != nil {
		return
	}

	// get ABI to parse
	file, err := os.Open(config.DefaultConfig.ABIPath)
	if err != nil {
		return
	}
	contractABI, err := abi.JSON(file)
	if err != nil {
		return
	}

	for _, record := range logs {
		res, err := parseEvt(contractABI, record.Data)
		if err != nil {
			continue
		}
		// do something to compare here
		fullResult = append(fullResult, res)
	}
	return fullResult, nil
}
