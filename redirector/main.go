package main

import (
	"context"
	"fmt"
	"gr-blockchain-side/internal"
	"gr-blockchain-side/internal/config"
	"log"

	"github.com/gin-gonic/gin"
)

var (
	engine *gin.Engine
	ctx    context.Context
	cancel context.CancelFunc
)

func init() {
	fmt.Println("Initializing.....")
	if err := internal.LoadEnv(); err != nil {
		log.Fatal(err)
	}

	if err := internal.LoadConfiguration(); err != nil {
		log.Fatal(err)
	}
	ctx, cancel = context.WithCancel(context.Background())
	engine = internal.InitEngine(config.DefaultConfig.DeployMode)
}

func main() {
	internal.SetUpAPIs(engine)
	internal.SetupGracefulShutdown(ctx, config.DefaultConfig.DefaultPort, engine)
	cancel()
}
