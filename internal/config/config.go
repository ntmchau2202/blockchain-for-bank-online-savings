package config

type Config struct {
	DefaultPort string
	Node        string // link to the blockchain remote node
	DeployMode  string // test or release
}

var DefaultConfig Config
