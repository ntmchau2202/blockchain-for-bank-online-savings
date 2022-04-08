package signer

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type Signer struct {
	passphrase string
	privateKey string
	opts       *bind.TransactOpts
}

func (s *Signer) GetOwner() (addr common.Address) {
	return s.opts.From
}

func (s *Signer) Sign(tx *types.Transaction) (signedTx *types.Transaction, err error) {
	return s.opts.Signer(s.GetOwner(), tx)
}

func (s *Signer) GetTxnOptsOfSigner() *bind.TransactOpts {
	return s.opts
}
