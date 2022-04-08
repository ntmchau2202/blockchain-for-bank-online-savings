package signer

import (
	"encoding/json"
	"io/ioutil"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
)

func NewSignerWithChainID(privateKeyPath string, chainId int64) (s *Signer, err error) {
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

	chainID := new(big.Int).SetInt64(chainId)
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		return
	}
	s = &Signer{}
	s.opts = auth

	return
}
