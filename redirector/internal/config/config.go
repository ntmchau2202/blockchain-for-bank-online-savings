package config

import "gr-blockchain-side/internal/blockchain/signer"

type Config struct {
	DefaultPort string
	Node        string // link to the blockchain remote node
	DeployMode  string // test or release
	Scanner     string // link to default inf scanner
	APIKey      string // api key to use the inf scanner

	Signer   *signer.Signer
	Contract string
	ABIPath  string
}

var DefaultConfig Config
