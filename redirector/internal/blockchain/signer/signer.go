package signer

import (
	"encoding/json"
	"io/ioutil"
	"math/big"
	"os"
	"strconv"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

type Signer struct {
	Keystore   string `json:"keystore_path"`
	Passphrase string `json:"passphrase"`
	opts       *bind.TransactOpts
}

func (signer *Signer) GetTransactOpts() *bind.TransactOpts {
	return signer.opts
}

func (signer *Signer) GetAddress() ethereum.Address {
	return signer.opts.From
}

func (signer *Signer) Sign(tx *types.Transaction) (*types.Transaction, error) {
	return signer.opts.Signer(signer.GetAddress(), tx)
}

func NewSigner(privateKeyPath string) (signer *Signer, err error) {
	// load private key file here
	raw, err := ioutil.ReadFile(privateKeyPath)
	if err != nil {
		return
	}
	var data map[string]string = make(map[string]string)
	err = json.Unmarshal(raw, &data)
	if err != nil {
		return
	}
	rawPrivateKey := data["key"]
	privateKey, err := crypto.HexToECDSA(rawPrivateKey)
	if err != nil {
		return
	}

	chainIDInt, err := strconv.ParseInt(os.Getenv("CHAIN_ID"), 10, 64)
	if err != nil {
		return
	}
	chainID := big.NewInt(chainIDInt)
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		return
	}
	signer = &Signer{}
	signer.opts = auth

	return
}
