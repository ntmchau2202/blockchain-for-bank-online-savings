package internal

import (
	"os"

	"blockchain-server/internal/config"

	"github.com/joho/godotenv"
)

func LoadEnv() (err error) {
	if err = godotenv.Load(os.Args[1]); err != nil {
		return
	}
	config.DefaultConfig.DefaultPort = os.Getenv("DEFAULT_PORT")
	config.DefaultConfig.DeployMode = os.Getenv("MODE")
	config.DefaultConfig.Node = os.Getenv("BLOCKCHAIN_NODE")
	return nil
}
